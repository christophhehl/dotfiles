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

func askAboutMissing(name string) (string, error) {
	greenArrow := color.GreenString("==>")
	fmt.Println(greenArrow, color.HiMagentaString(name), "is not installed. Do you want to install it now?")
	fmt.Println(greenArrow, color.CyanString("[Y]es"), "[N]o [Ne]ver")
	fmt.Print(greenArrow, " ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	// remove trailing newline character.
	input = input[:len(input)-1]

	if strings.EqualFold(input, "y") || input == "" {
		return "y", nil
	} else if strings.EqualFold(input, "n") {
		return "n", nil
	} else if strings.EqualFold(input, "ne") {
		return "ne", nil
	} else {
		return "", fmt.Errorf("invalid option")
	}
}

func handleMissing(name string) bool {
	answer, err := askAboutMissing(name)
	if err != nil {
		fmt.Println("Invalid option. Aborting...")
		os.Exit(1)
	}
	switch answer {
	case "y":
		fmt.Println(":: Installing " + name + "...")
		return installProgram(name)
	case "n":
		fmt.Println("Skipping install. This will be asked again next time.")
		return false
	case "ne":
		fmt.Println("Skipping install. This will", color.RedString("not"), "be asked again next time.")
		saveIgnoreMissing(name)
		return false
	default:
		return false
	}
}

func saveIgnoreMissing(name string) {
	panic("unimplemented")
}

func installProgram(name string) bool {
	if name == "yay" {
		return installYay()
	} else if slices.Contains(yayInstallable, name) {
		return installProgramViaYay(name)
	} else {
		return installSpecial(name)
	}
}

func installYay() bool {
	cmd := exec.Command("sudo", "ls") // TODO: change. only placeholder for sudo execution code

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("sudo command failed: %s, output: %s", err, output)
	}

	fmt.Printf("sudo command output: %s\n", output)
	panic("unimplemented")
}

func installProgramViaYay(name string) bool {
	panic("unimplemented")
}

func installSpecial(name string) bool {
	panic("unimplemented")
}

var yayInstallable = []string{"neovim"}

func main() {
	if os.Geteuid() == 0 {
		fmt.Println("Avoid running yay as root/sudo.")
		os.Exit(1)
	}

	// schauen installed, falls nicht installieren? [ja|ignore|ignore_always(save to ~/.dotfiles.cfg)]
	fmt.Println(":: Checking installed programs...")

	isZshInstalled := checkInstalled("zsh")
	isOmzInstalled := checkInstalled("oh-my-zsh")
	isP10kInstalled := checkInstalled("powerlevel10k")
	isYayInstalled := checkInstalled("yay")

	fmt.Println()

	if !isZshInstalled {
		handleMissing("zsh")
	}
	if !isOmzInstalled {
		handleMissing("oh-my-zsh")
	}
	if !isP10kInstalled {
		handleMissing("powerlevel10k")
	}
	if !isYayInstalled {
		isYayInstalled = handleMissing("yay")
	}
}
