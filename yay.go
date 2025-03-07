package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func checkYayInstalled() bool {
	_, err := exec.LookPath("yay")
	isInstalled := err == nil

	if isInstalled {
		fmt.Println("- yay", strings.Repeat(" ", 20-len("yay")), color.GreenString("(installed)"))
	} else {
		fmt.Println("- yay", strings.Repeat(" ", 20-len("yay")), color.RedString("(not installed)"))
	}
	fmt.Println()

	return isInstalled
}

func handleYayMissing() bool {
	fmt.Println(greenArrow+color.MagentaString("yay"), "is not installed. Do you want to install it now?")

	answer := getYesNoAnswer()
	if answer == "yes" {
		fmt.Println(":: Installing yay...")
		return installProgram("yay")
	} else {
		fmt.Println("Skipping install. This will be asked again next time.")
	}
	return false
}
