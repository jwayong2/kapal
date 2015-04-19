
package sysfs

import(
	"os"
	"bufio"
	"log"
	"errors"
	"strings"
	"fmt"
	"os/exec"
	"path/filepath"
)


// ListDevices prints a list of all storage devices attached to a linux machine
func ListDevicesCmd() {
	for _, dev := range findDevices() {
		fmt.Printf("--> %s", dev)
	}
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

func getDeviceName(line string) (string, error){
	parts := strings.Fields(line)
	if len(parts) >= 4 && isStorageDevice(parts[3]){
		return parts[3], nil
	}
	return "", errors.New("Line does not contain any storage device")
}

func isStorageDevice(devName string) bool {
	return strings.HasPrefix(devName, "hd") || strings.HasPrefix(devName, "sd")
}

func recognizeTypeDevice(devName string) string {
	output, err := exec.Command("file", "-s", filepath.Join("/dev", devName)).CombinedOutput()
	if err != nil{
		log.Fatal(err)
	}
	return string(output)

}
