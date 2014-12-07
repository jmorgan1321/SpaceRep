package main

import (
	"fmt"
	"os/exec"
)

// TODO: make version for different builds

func openBrowser(url string) {
	cmd := exec.Command("cmd", "/c", "start", url)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%v\n", url)
	}
}
