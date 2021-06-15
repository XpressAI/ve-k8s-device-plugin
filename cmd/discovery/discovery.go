package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

const filePath = `C:\Users\Hazim.Hasnan\Documents\Work\vedummy.txt`

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
		return
	}

	vecmdResult := strings.NewReader(string(stdout))
	scanner := bufio.NewScanner(vecmdResult)
	for scanner.Scan() {
		trimmedLine := strings.TrimSpace(scanner.Text())
		lines = append(lines, trimmedLine)
	}

	fmt.Println(lines[0])

	// f, err := os.Open(filePath)
	// if err != nil {
	// 	return
	// }

	// scanner := bufio.NewScanner(f)
	// for scanner.Scan() {
	// 	trimmedLine := strings.TrimSpace(scanner.Text())
	// 	lines = append(lines, trimmedLine)
	// }
	// f.Close()

	for _, line := range lines {

		switch {
		case strings.HasPrefix(line, "Attached VEs"):

			arrayString := strings.Fields(line)
			numString := arrayString[len(arrayString)-1]
			num, err := strconv.ParseInt(numString, 0, 32)
			if err != nil {
				fmt.Println("Unable to parse string to int")
				fmt.Println(err)
				return
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
			currentVe["dev"] = "/dev/ve" + string(veID-1)
			stat := unix.Stat_t{}
			err := unix.Stat("/dev/ve", &stat)
			if err != nil {
				fmt.Println("Error getting major and minor number")
			}
			major := unix.Major(stat.Rdev)
			minor := unix.Minor(stat.Rdev)
			fmt.Println(major, minor)
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

	for key, value := range ves {
		fmt.Println("%v : %v", key, value)
	}
}
