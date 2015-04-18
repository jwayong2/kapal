package sysfs




// ListDevices prints a list of all storage devices attached to a linux machine
func ListDevicesCmd() {
	for _, dev := range findDevices() {
		fmt.Printf("--> %s\n", dev)
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
		if err != nil {
			devices = append(devices, device)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return devices
}

func getDeviceName(line string) (string, error){
	parts, err := line.Split()
	if len(parts) == 4 && isStorageDevice(parts[3]){
		return parts[3]
	}
	return errors.New("Line does not contain any storage device")
}

func isStorageDevice(devName string) bool {
	return strings.HasPrefix(devName, "hd") || strings.HasPrefix(devName, "sd")
}
