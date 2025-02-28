package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func checkInstalled(name string) bool {
	var isInstalled bool // Declare the result variable outside the switch

	switch name {
	case "zsh":
		isInstalled = checkZsh()
	case "yay":
		isInstalled = checkYay()
	case "oh-my-zsh":
		isInstalled = checkOmz()
	case "powerlevel10k":
		isInstalled = checkP10k()
	default:
		fmt.Println("-", name, strings.Repeat(" ", 20-len(name)), color.YellowString("(unknown check)"))
		return false
	}

	if isInstalled {
		fmt.Println("-", name, strings.Repeat(" ", 20-len(name)), color.GreenString("(installed)"))
	} else {
		fmt.Println("-", name, strings.Repeat(" ", 20-len(name)), color.RedString("(not installed)"))
	}

	return isInstalled
}

func checkZsh() bool {
	_, err := exec.LookPath("zsh")
	return err == nil
}

func checkYay() bool {
	_, err := exec.LookPath("yay")
	return err == nil
}

func checkOmz() bool {
	homeDir := getHomeDir()
	ohMyZshPath := homeDir + "/.oh-my-zsh"

	_, err := os.Stat(ohMyZshPath)
	return !os.IsNotExist(err)
}

func checkP10k() bool {
	homeDir := getHomeDir()
	zshCustomDir := os.Getenv("ZSH_CUSTOM")
	if zshCustomDir == "" {
		zshCustomDir = filepath.Join(homeDir, ".oh-my-zsh", "custom")
	}

	powerlevel10kDir := filepath.Join(zshCustomDir, "themes", "powerlevel10k")
	_, err := os.Stat(powerlevel10kDir)
	return !os.IsNotExist(err)
}

func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}
	return homeDir
}

func main() {
	// schauen installed, falls nicht installieren? [ja|ignore|ignore_always(save to .dotfiles.cfg)]
	fmt.Println(":: Checking installed programs...")

	isZshInstalled := checkInstalled("zsh")
	isOmzInstalled := checkInstalled("oh-my-zsh")
	isP10kInstalled := checkInstalled("powerlevel10k")
	isYayInstalled := checkInstalled("yay")

	fmt.Println()

	if isZshInstalled {
	}
	if isOmzInstalled {
	}
	if isP10kInstalled {
	}
	if !isYayInstalled {
		greenArrow := color.GreenString("==>")
		fmt.Println(greenArrow, "Yay is not installed. Do you want to install it now?")
		fmt.Println(greenArrow, "[Y]es [N]o [Ne]ver")
		fmt.Print(greenArrow, " ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')

		// remove trailing newline character.
		input = input[:len(input)-1]

		if strings.EqualFold(input, "y") {
			fmt.Print("Installing... beep boop")
		} else if strings.EqualFold(input, "n") {
			fmt.Println("Ok, doing nothing.")
		} else if strings.EqualFold(input, "ne") {
			fmt.Print("Info is being stored. You wont be asked again")
		}
	}
}
