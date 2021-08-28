package main

import (
	"github.com/golang/glog"
	"io/ioutil"
	"os"
	"strings"
)

type VEInfo struct {
	slot   string
	device string
}

func GetNECVEInfo() ([]VEInfo, error) {
	var ves = make([]VEInfo, 0)

	devices, err := ioutil.ReadDir("/dev")
	if err != nil {
		glog.Error(err, "Error listing /dev")
		return nil, err
	}

	foundDevicesCount := 0
	foundDevicesOnline := 0
	for _, deviceFile := range devices {
		deviceName := deviceFile.Name()
		if strings.HasPrefix(deviceName, "veslot") {
			glog.Info("VE Device --- ", deviceName)
			foundDevicesCount++

			device, err := os.Readlink("/dev/" + deviceName)
			if err != nil {
				glog.Error(err, "Error reading "+deviceName+" symlink")
				return nil, err
			}
			glog.Info("VE actual device: ", device)

			bytes, err := ioutil.ReadFile("/sys/class/ve/" + device + "/ve_state")
			if err != nil {
				glog.Error(err, "Error reading ve_state")
				return nil, err
			}

			state := "OFFLINE"
			if strings.TrimSpace(string(bytes)) == "1" {
				state = "ONLINE"
			}
			glog.Info("VE State: ", state)

			if state == "ONLINE" {
				foundDevicesOnline++
				var curVe = make(map[string]interface{})
				curVe["slot"] = deviceName
				curVe["device"] = device

				ves = append(ves, VEInfo{
					slot:   deviceName,
					device: device,
				})
			}
		}
	}
	glog.Info("Total Devices found:  ", foundDevicesCount)
	glog.Info("ONLINE Devices found:  ", foundDevicesOnline)
	return ves, nil
}
