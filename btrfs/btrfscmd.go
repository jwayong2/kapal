package btrfscmd

import (
	"fmt"
	"os/exec"
	"path"
	"strings"
	"bytes"
)

func IsBtrfsVolume(volume string) error {
	cmd := exec.Command("btrfs", "filesystem", "df", volume)
	return cmd.Run()
}

func MakeDeviceBtrfs(device string) error {
	cmd := exec.Command("mkfs.btrfs", "-f", device)
	return cmd.Run()
}

func SubvolumeCreate (volume string, name string) error {
	cmd := exec.Command("btrfs", "subvolume", "create", path.Join(volume, name))
	return cmd.Run()
}

func SubvolumeList (volume string) (result []string) {
	cmd := exec.Command("btrfs", "subvolume", "list", "-a", volume)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	
	if err != nil {
		fmt.Errorf("Failed to list subvolumes on root volume: %s",volume)
		return
	} else {
		lines := strings.Split(out.String(), "\n")
		for _, line := range lines {
			if tokens := strings.Split(line, "path"); len(tokens) != 2 {
				fmt.Errorf("Can't parse subvolume on line: %s", line)
			} else {
				result = append(result, strings.TrimSpace(tokens[1]))
			}
		}
	}
	return
}	

func SubvolumeSnapshot (volume string, sourcesubvolume string, targetsubvolume string, readonly bool) error {
	var cmd *exec.Cmd
	if readonly {
		cmd = exec.Command("btrfs", "subvolume", "snapshot" ,"-r", path.Join(volume,sourcesubvolume), path.Join(volume,targetsubvolume))
	} else {
		cmd = exec.Command("btrfs", "subvolume", "snapshot", path.Join(volume,sourcesubvolume), path.Join(volume,targetsubvolume))
	}
	return cmd.Run()
}

func SubvolumeDelete (volume string, name string) error {
	cmd := exec.Command("btrfs", "subvolume", "delete", path.Join(volume, name))
	return cmd.Run()
}

func Sync() error {
	cmd := exec.Command("sync")
	return cmd.Run()
}

/* TODO: BTRFS Send and Recv */
