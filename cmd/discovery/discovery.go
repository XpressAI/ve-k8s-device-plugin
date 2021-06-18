package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func main() {

	var veCount int
	var veID int
	var currentVe = make(map[string]interface{})
	var ves = make([]map[string]interface{}, 0)
	var lines []string

	//Replace the arguments to "vecmd", "info"
	cmd := exec.Command("cat", "C:\\Users\\Hazim.Hasnan\\Documents\\Work\\vedummy.txt")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println("No Result to from the Command")
		fmt.Println(err.Error())
	}

	//trim and create array of lines
	vecmdResult := strings.NewReader(string(stdout))
	scanner := bufio.NewScanner(vecmdResult)
	for scanner.Scan() {
		trimmedLine := strings.TrimSpace(scanner.Text())
		lines = append(lines, trimmedLine)
	}

	for _, line := range lines {

		switch {
		case strings.HasPrefix(line, "Attached VEs"):

			arrayString := strings.Fields(line)
			numString := arrayString[len(arrayString)-1]
			num, err := strconv.ParseInt(numString, 0, 32)
			if err != nil {
				fmt.Println("Unable to parse string to int")
				fmt.Println(err)

			}
			veCount = int(num)
			fmt.Println("Number of Attached VEs --- ", veCount)

		case strings.HasPrefix(line, "[V"):
			fmt.Println(line)
			if len(currentVe) != 0 {
				ves = append(ves, currentVe)
			}
			veID++
			currentVe["id"] = veID
			currentVe["dev"] = "/dev/ve" + strconv.FormatInt(int64(veID-1), 10)
			// stat := unix.Stat_t{}
			// err := unix.Stat("/dev/ve", &stat)
			// if err != nil {
			// 	fmt.Println("Error getting major and minor number")
			// }
			// major := unix.Major(stat.Rdev)
			// minor := unix.Minor(stat.Rdev)
			// fmt.Println(major, minor)
			//fmt.Println("Major:", uint64(stat.Rdev/256), "Minor:", uint64(stat.Rdev%256))

		case strings.HasPrefix(line, "VE State"):

			arrayString := strings.Fields(line)
			state := arrayString[len(arrayString)-1]
			currentVe["state"] = state
			fmt.Println("VE State is --- ", currentVe["state"])

		case strings.HasPrefix(line, "Bus ID"):

			arrayString := strings.Fields(line)
			busId := arrayString[len(arrayString)-1]
			currentVe["busId"] = busId
			fmt.Println("Bus ID is --- ", currentVe["busId"])
		}

	}
	ves = append(ves, currentVe)

	for key, value := range ves {
		fmt.Println("---------------")
		for key, value := range value {
			fmt.Printf("%v : %v\n", key, value)
		}
		fmt.Printf("Last Line: %v : %v\n", key, value)
	}

	y := getVEs(ves)
	fmt.Println(y)
}

//Extract the id and dev from NECVE metadatas
func getVEs(NECDevs []map[string]interface{}) map[string]string {
	ves := make(map[string]string)

	for _, values := range NECDevs {
		id := fmt.Sprint(values["id"])
		ves[id] = fmt.Sprint(values["dev"])
	}
	return ves
}
