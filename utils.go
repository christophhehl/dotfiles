package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/fatih/color"
)

var greenArrow = color.GreenString("==>") + " "

func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}
	return homeDir
}

func getInput() string {
	fmt.Print(greenArrow)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	// remove trailing newline character.
	return input[:len(input)-1]
}

func getYesNoAnswer() string {
	fmt.Println(greenArrow+color.CyanString("[Y]es"), "[N]o")

	input := getInput()

	answer := "no"
	if strings.EqualFold(input, "y") || strings.EqualFold(input, "yes") || input == "" {
		answer = "yes"
	}
	return answer
}

func executeCommandDir(sudo bool, command string, args []string, dir string) bool {
	var cmd *exec.Cmd
	if sudo {
		sudoArgs := append([]string{command}, args...)
		cmd = exec.Command("sudo", sudoArgs...)
	} else {
		cmd = exec.Command(command, args...)
	}

	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println(string(output))
	return true
}

func installProgram(name string) bool {
	if slices.Contains(yayInstallable, name) {
		return installProgramViaYay(name)
	} else {
		return installSpecial(name)
	}
}

func installProgramViaYay(name string) bool {
	panic("unimplemented")
}

func installSpecial(name string) bool {
	switch name {
	case "yay":
		if successful := executeCommandDir(true, "pacman", []string{"-S", "--needed", "git", "base-devel"}, ""); !successful {
			fmt.Println("todo")
			return false
		}
		if successful := executeCommandDir(false, "git", []string{"clone", "https://aur.archlinux.org/yay.git", "yay_install"}, getHomeDir()); !successful {
			fmt.Println("todo")
			return false
		}
		defer func() { // ensure cleanup
			err := os.RemoveAll(filepath.Join(getHomeDir(), "yay_install"))
			if err != nil {
				fmt.Println("Couldn't delete folder needed for installation of yay. Check your home directory for \"yay_install\" folder, if you want to remove it yourself.")
			}
		}()
		if successful := executeCommandDir(false, "makepkg", []string{"-si"}, filepath.Join(getHomeDir(), "yay_install")); !successful {
			fmt.Println("todo")
			return false
		} else {
			return successful
		}
	}
	return false
}
