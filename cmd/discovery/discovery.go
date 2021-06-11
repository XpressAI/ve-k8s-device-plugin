package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const filePath = `C:\Users\Hazim.Hasnan\Documents\Work\vedummy.txt`

func main() {

	var veCount int
	//var veID int
	var currentVe = make(map[string]interface{})
	//  var ves = make([]map[string]interface{}, 0, 5)

	var lines []string

	// cmd := exec.Command("go", "version")
	// stdout, err := cmd.Output()

	// if err != nil {
	// 	fmt.Println("No Result to from the Command")
	// 	fmt.Println(err.Error())
	// 	return
	// }
	//vecmdResult := strings.NewReader(string(stdout))
	// for scanner.Scan() {
	// 	lines = append(lines, scanner.Text())
	// }

	// fmt.Println(lines[0])
	// scanner := bufio.NewScanner(vecmdResult)

	f, err := os.Open(filePath)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	f.Close()

	for _, line := range lines {

		switch {
		case strings.HasPrefix(line, "Attached VEs"):

			fmt.Println(line)
			content := strings.Fields(line)
			numString := content[len(content)-1]
			num, err := strconv.ParseInt(numString, 0, 32)
			if err != nil {
				fmt.Println("Unable to parse string to int")
				fmt.Println(err)
				return
			}
			veCount = int(num)
			fmt.Println(veCount)

		case strings.HasPrefix(line, "[V"):
			fmt.Println(line)

		case strings.HasPrefix(line, "VE State"):

			fmt.Println(line)
			content := strings.Fields(line)
			numString := content[len(content)-1]
			num, err := strconv.ParseInt(numString, 0, 32)
			if err != nil {
				fmt.Println("Unable to parse string to int")
				fmt.Println(err)
				return
			}
			currentVe["state"] = num

		case strings.HasPrefix(line, "Bus ID"):
			fmt.Println(line)
			content := strings.Fields(line)
			numString := content[len(content)-1]
			num, err := strconv.ParseInt(numString, 0, 32)
			if err != nil {
				fmt.Println("Unable to parse string to int")
				fmt.Println(err)
				return
			}
			currentVe["busId"] = num
		}

	}

	cmd := exec.Command("go", "version")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println("No Result to from the Command")
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(stdout))

}
