package sysfs

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

// ListDevices prints a list of all storage devices attached to a linux machine
func ListDevicesCmd() {
	fmt.Printf("Current storage devices:\n\n")
	for _, dev := range findDevices() {
		fmt.Printf("%s\n", dev)
	}
	if !isSudoerUser() {
		fmt.Printf("\nPlease, rerun this command with sudo if you want to learn more information about these devices")
	}
	fmt.Printf("\n")
}

func findDevices() []string {
	var devices []string

	file, err := os.Open("/proc/partitions")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		device, err := getDeviceName(scanner.Text())
		if err == nil {
			devDesc := recognizeTypeDevice(device)
			devices = append(devices, devDesc)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return devices
}

func getDeviceName(line string) (string, error) {
	parts := strings.Fields(line)
	if len(parts) >= 4 && isStorageDevice(parts[3]) {
		return parts[3], nil
	}
	return "", errors.New("Line does not contain any storage device")
}

func isStorageDevice(devName string) bool {
	return strings.HasPrefix(devName, "hd") || strings.HasPrefix(devName, "sd")
}

func recognizeTypeDevice(devName string) string {
	if isSudoerUser() {
		output, err := exec.Command("file", "-s", filepath.Join("/dev", devName)).CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}
		return strings.TrimRight(string(output), "\n")
	}
	return "/dev/" + devName
}

func isSudoerUser() bool {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	// If not superuser
	if user.Uid != "0" {
		err := exec.Command("sudo", "-n", "btrfs", "help").Run()
		return err == nil
	}

	return true
}
