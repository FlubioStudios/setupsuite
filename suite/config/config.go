package config

import (
	"fmt"
	"io/ioutil"
	"strings"
	"suite/suite/utils"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ReadConfig() {
	fmt.Println("I read")
	dat, err := ioutil.ReadFile("/etc/setupsuite/config.sscfg")
	check(err)
	fmt.Println("//////////////////////!")
	fmt.Print(string(dat))
	fmt.Println("//////////////////////!")

	//split config into lines
	configlines := strings.Fields(string(dat))

	//reverse configlines
	for i, j := 0, len(configlines)-1; i < j; i, j = i+1, j-1 {
		configlines[i], configlines[j] = configlines[j], configlines[i]
	}

	length := len(configlines)

	for length > 0 {
		line := configlines[length-1]
		fmt.Println("Index: ", length)
		//check if it's a config-group
		if strings.HasPrefix(line, ".") {
			if strings.HasSuffix(line, "{") {
				//replace config group indicators
				replaced := strings.Replace(line, ".", "", -1)
				//find ending }
				for x := len(configlines) - 1; x < utils.FindIndex(configlines, line, 0) && x > 0; x-- {
					if configlines[x] == "}" {
						fmt.Println("line below: ", line, ": ", configlines[x])
					}
				}
				fmt.Println("Start of config group: ", strings.Replace(replaced, "{", "", -1))
			}
		}
		fmt.Println(configlines[length-1])
		length--
	}

}
