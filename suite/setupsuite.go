package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	if runtime.GOOS == "windows" {
		fmt.Println("The LINUX setup suite is only to be used on linux based systems")
	} else {
		fmt.Println("Starting Serversetup...")
		fmt.Println(len(os.Args), os.Args)
		execute()
	}
}

func execute() {

	// let's try the whoami command here
	out, err := exec.Command("whoami").Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	output := string(out[:])
	if output != "root" {
		fmt.Println("The Setupsuite can only be run as root")
	}

	// here we perform the pwd command.
	// we can store the output of this in our out variable
	// and catch any errors in err
	out, err = exec.Command("ls").Output()

	// if there is an error with our execution
	// handle it here
	if err != nil {
		fmt.Printf("%s", err)
	}

	// as the out variable defined above is of type []byte we need to convert
	// this to a string or else we will see garbage printed out in our console
	// this is how we convert it to a string
	fmt.Println("\"ls\" Successfully Executed")
	output = string(out[:])
	fmt.Println(output)

}
