package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
	"golang.org/x/net/context"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

// Plugin is identical to DevicePluginServer
// interface of device plugin API.
type Plugin struct {
	NECVEs    map[string]string
	Heartbeat chan bool
}

// Start is an optional interface that could be implemented by plugin.
// If case Start is implemented, it will be executed by Manager after
// plugin instantiation and before its registration to kubelet. This
// method could be used to prepare resources before they are offered
// to Kubernetes.
func (p *Plugin) Start() error {
	return nil
}

// Stop is an optional interface that could be implemented by plugin.
// If case Stop is implemented, it will be executed by Manager after the
// plugin is unregistered from kubelet. This method could be used to tear
// down resources.
func (p *Plugin) Stop() error {
	return nil
}

func simpleHealthCheck() bool {
	var kfd *os.File
	var err error
	if kfd, err = os.Open("/dev/kfd"); err != nil {
		glog.Error("Error opening /dev/kfd")
		return false
	}
	kfd.Close()
	return true
}

// GetDevicePluginOptions returns options to be communicated with Device
// Manager
func (p *Plugin) GetDevicePluginOptions(ctx context.Context, e *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{}, nil
}

// PreStartContainer is expected to be called before each container start if indicated by plugin during registration phase.
// PreStartContainer allows kubelet to pass reinitialized devices to containers.
// PreStartContainer allows Device Plugin to run device specific operations on the Devices requested
func (p *Plugin) PreStartContainer(ctx context.Context, r *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return &pluginapi.PreStartContainerResponse{}, nil
}

// ListAndWatch returns a stream of List of Devices
// Whenever a Device state change or a Device disappears, ListAndWatch
// returns the new list
func (p *Plugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {

	necDevs, err := GetNECVEInfo()
	if err != nil {
		glog.Error(err, "Error reading veinfo")
	}
	p.NECVEs = getVEs(necDevs)
	glog.Info(p.NECVEs)
	devs := make([]*pluginapi.Device, len(p.NECVEs))

	func() {
		i := 0
		for id := range p.NECVEs {
			// idStr := strconv.FormatInt(int64(id), 10)
			dev := &pluginapi.Device{
				ID:     id,
				Health: pluginapi.Healthy,
			}
			devs[i] = dev
			i++

			// numas, err := hw.GetNUMANodes(id)
			// glog.Infof("Watching GPU with bus ID: %s NUMA Node: %+v", id, numas)
			// if err != nil {
			// 	glog.Error(err)
			// 	continue
			// }

			// if len(numas) == 0 {
			// 	glog.Errorf("No NUMA for GPU ID: %s", id)
			// 	continue
			// }

			// numaNodes := make([]*pluginapi.NUMANode, len(numas))
			// for j, v := range numas {
			// 	numaNodes[j] = &pluginapi.NUMANode{
			// 		ID: int64(v),
			// 	}
			// }
			////testingbranch
			// dev.Topology = &pluginapi.TopologyInfo{
			// 	Nodes: numaNodes,
			// }
		}
	}()

	s.Send(&pluginapi.ListAndWatchResponse{Devices: devs})

	for { //****                                                                       *****WHy no error****

		<-p.Heartbeat
		var health = pluginapi.Unhealthy

		// TODO there are no per device health check currently
		// TODO all devices on a node is used together by kfd
		if simpleHealthCheck() {
			health = pluginapi.Healthy
		}

		for i := 0; i < len(p.NECVEs); i++ {
			devs[i].Health = health
		}
		s.Send(&pluginapi.ListAndWatchResponse{Devices: devs})
	}
	// returning a value with this function will unregister the plugin from k8s
}

// GetPreferredAllocation returns a preferred set of devices to allocate
// from a list of available ones. The resulting preferred allocation is not
// guaranteed to be the allocation ultimately performed by the
// devicemanager. It is only designed to help the devicemanager make a more
// informed allocation decision when possible.
func (p *Plugin) GetPreferredAllocation(ctx context.Context, r *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	response := pluginapi.PreferredAllocationResponse{}
	return &response, nil
}

// Allocate is called during container creation so that the Device
// Plugin can run device specific operations and instruct Kubelet
// of the steps to make the Device available in the container
func (p *Plugin) Allocate(ctx context.Context, r *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	var response pluginapi.AllocateResponse
	var car pluginapi.ContainerAllocateResponse
	var dev *pluginapi.DeviceSpec

	for _, req := range r.ContainerRequests {
		car = pluginapi.ContainerAllocateResponse{}

		for _, id := range req.DevicesIDs {
			glog.Infof("Allocating device ID: %s", id)

			for _, v := range p.NECVEs { //initially p.AMDGPUs[id]
				devpath := v //"dev/ve0"
				dev = new(pluginapi.DeviceSpec)
				dev.HostPath = devpath
				dev.ContainerPath = devpath
				dev.Permissions = "rw"
				car.Devices = append(car.Devices, dev)
			}
		}

		response.ContainerResponses = append(response.ContainerResponses, &car)
	}

	return &response, nil
}

// Lister serves as an interface between imlementation and Manager machinery. User passes
// implementation of this interface to NewManager function. Manager will use it to obtain resource
// namespace, monitor available resources and instantate a new plugin for them.
type Lister struct {
	ResUpdateChan chan dpm.PluginNameList
	Heartbeat     chan bool
}

// GetResourceNamespace must return namespace (vendor ID) of implemented Lister. e.g. for
// resources in format "color.example.com/<color>" that would be "color.example.com".
func (l *Lister) GetResourceNamespace() string {
	return "nec.com"
}

// Discover notifies manager with a list of currently available resources in its namespace.
// e.g. if "color.example.com/red" and "color.example.com/blue" are available in the system,
// it would pass PluginNameList{"red", "blue"} to given channel. In case list of
// resources is static, it would use the channel only once and then return. In case the list is
// dynamic, it could block and pass a new list each times resources changed. If blocking is
// used, it should check whether the channel is closed, i.e. Discover should stop.
func (l *Lister) Discover(pluginListCh chan dpm.PluginNameList) {
	for {
		select {
		case newResourcesList := <-l.ResUpdateChan: // New resources found
			pluginListCh <- newResourcesList
		case <-pluginListCh: // Stop message received
			// Stop resourceUpdateCh
			return
		}
	}
}

// NewPlugin instantiates a plugin implementation. It is given the last name of the resource,
// e.g. for resource name "color.example.com/red" that would be "red". It must return valid
// implementation of a PluginInterface.
func (l *Lister) NewPlugin(resourceLastName string) dpm.PluginInterface {
	return &Plugin{
		Heartbeat: l.Heartbeat,
	}
}

//var gitDescribe string

func main() {
	versions := [...]string{
		"NEC device plugin for Kubernetes",
		"Vector 1",
	}

	flag.Usage = func() {
		for _, v := range versions {
			fmt.Fprintf(os.Stderr, "%s\n", v)
		}
		fmt.Fprintln(os.Stderr, "Usage:")
		flag.PrintDefaults()
	}
	var pulse int
	flag.IntVar(&pulse, "pulse", 0, "time between health check polling in seconds.  Set to 0 to disable.")
	flag.Parse()

	for _, v := range versions {
		glog.Infof("%s", v)
	}

	l := Lister{
		ResUpdateChan: make(chan dpm.PluginNameList),
		Heartbeat:     make(chan bool),
	}
	manager := dpm.NewManager(&l)

	if pulse > 0 {
		go func() {
			glog.Infof("Heart beating every %d seconds", pulse)
			for {
				time.Sleep(time.Second * time.Duration(pulse))
				l.Heartbeat <- true
			}
		}()
	}

	go func() {
		var path = "/sys/class/kfd"
		if _, err := os.Stat(path); err == nil {
			l.ResUpdateChan <- []string{"ve"} //Change to VE
		}
	}()
	manager.Run()

}
