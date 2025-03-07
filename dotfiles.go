package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var yayInstallable = []string{"neovim"}

func installPreset() string {
	// Give install selection
	fmt.Println(":: Available presets...")
	fmt.Println("  " + color.MagentaString("3") + "  server")
	fmt.Println("  " + color.MagentaString("2") + "  laptop")
	fmt.Println("  " + color.MagentaString("1") + "  desktop")
	fmt.Println(greenArrow + "Select installation preset: (eg: \"1\" or \"laptop\")")
	selection := getInput()

	switch trimmedSelection := strings.ToLower(strings.TrimSpace(selection)); trimmedSelection {
	case "3", "server":
		fmt.Println(":: Installing server packages...")
		if success := installPresetPackages("server"); !success {
			return "error"
		}
		return "server"
	case "2", "laptop":
		fmt.Println(":: Installing laptop packages...")
		if success := installPresetPackages("laptop"); !success {
			return "error"
		}
		return "laptop"
	case "1", "desktop":
		fmt.Println(":: Installing desktop packages...")
		if success := installPresetPackages("desktop"); !success {
			return "error"
		}
		return "desktop"
	default:
		fmt.Println(":: Skipping installation...")
		return ""
	}
}

func installPresetPackages(preset string) bool {
	for _, packet := range presetsPackages[preset] {
		success := installProgram(packet)
		if !success {
			return false
		}
	}
	return true
}

func copyPresetConfigs(preset string) bool {
	// TODO: copy files
	return false
}

var presetsPackages map[string][]string

func main() {
	if os.Geteuid() == 0 {
		fmt.Println("Avoid running dotfiles as root/sudo.")
		os.Exit(1)
	}

	presetsPackages = make(map[string][]string)
	presetsPackages["server"] = []string{""}
	presetsPackages["laptop"] = []string{""}
	presetsPackages["desktop"] = []string{""}

	searchValue := ""
	if len(os.Args) > 1 {
		searchValue := os.Args[1]
		searchValue += ""
	}

	if searchValue == "" {
		// Check for yay
		fmt.Println(":: Checking, if yay is installed...")
		isYayInstalled := checkYayInstalled()
		if !isYayInstalled {
			isYayInstalled = handleYayMissing()
		}

		if !isYayInstalled {
			fmt.Println(":: Yay is not installed, skipping installation...")
		} else {
			installedPreset := installPreset()
			if installedPreset == "error" {
				fmt.Println(":: Preset couldn't be installed.")
				fmt.Println(":: Exiting...")
				os.Exit(1)
			} else if installedPreset != "" {
				fmt.Println(greenArrow + "Installed preset: " + color.MagentaString(installedPreset))
				fmt.Println(greenArrow + "Would you like to copy config files for installed preset to your system?")
				answer := getYesNoAnswer()
				if answer == "no" {
					fmt.Println(":: No config files to copy...")
					fmt.Println(":: Exiting...")
					os.Exit(0)
				} else {
					success := copyPresetConfigs(installedPreset)
					if success {
						fmt.Println(":: Config files " + color.GreenString("successfully") + " installed.")
					} else {
						fmt.Println(":: Could " + color.RedString("not") + " copy config files.")
					}
					fmt.Println(":: Exiting...")
					os.Exit(0)
				}
			} else {
				fmt.Println(greenArrow + "Config files to copy: (eg: \"3\" or \"nvim\")")
				// TODO: choose config files to copy
			}
		}
	}
}
