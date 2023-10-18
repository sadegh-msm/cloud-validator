package main

import (
	"fmt"
	"os/exec"
)

func main() {
	output, err := exec.Command("pyhton3").Output()
	if err != nil {
		return
	}
	fmt.Println(string(output))
}
