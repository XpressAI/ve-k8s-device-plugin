package main

import (
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//Return slice of ves which is a map of each ve metadata
func GetNECVEInfo() ([]map[string]interface{}, error) {
	var ves = make([]map[string]interface{}, 0)

	files, err := ioutil.ReadDir("/sys/class/ve")

	if err != nil {
		glog.Error(err, "Error listing /sys/class/ve")
		return nil, err
	}

	glog.Info("Number of Attached VEs --- ", len(files))

	for _, file := range files {
		var curVe = make(map[string]interface{})
		devName := file.Name()
		glog.Info("VE Device --- ", devName)

		bytes, err := ioutil.ReadFile("/sys/class/ve/" + devName + "/ve_state")
		if err != nil {
			glog.Error(err, "Error reading ve_state")
			return nil, err
		}

		state := "OFFLINE"
		if strings.TrimSpace(string(bytes)) == "1" {
			state = "ONLINE"
		}
		glog.Info("VE State is --- ", state)

		readlink, err := os.Readlink("/sys/class/ve/" + devName + "/device")
		if err != nil {
			glog.Error(err, "Error reading device")
			return nil, err
		}

		fullBusId := filepath.Base(readlink)
		busIdParts := strings.SplitN(fullBusId, ":", 2)
		busId := busIdParts[1]
		glog.Info("Bus ID is --- ", busId)

		devNum := strings.Split(devName, "ve")[1]

		curVe["id"] = devNum
		curVe["dev"] = "/dev/" + devName
		curVe["state"] = state
		curVe["busId"] = busId
		ves = append(ves, curVe)
	}

	return ves, nil
}

func getVEs(NECDevs []map[string]interface{}) map[string]string {
	ves := make(map[string]string)

	for _, values := range NECDevs {
		id := fmt.Sprint(values["id"])
		ves[id] = fmt.Sprint(values["dev"])
	}
	return ves
}
