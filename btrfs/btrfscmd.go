package btrfscmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

func IsBtrfsVolume(volume string) error {
	cmd := exec.Command("btrfs", "filesystem", "df", volume)
	return cmd.Run()
}

func MakeDeviceBtrfs(device string) error {
	cmd := exec.Command("mkfs.btrfs", "-f", device)
	return cmd.Run()
}

func SubvolumeCreate(volume string, name string) error {
	cmd := exec.Command("btrfs", "subvolume", "create", path.Join(volume, name))
	return cmd.Run()
}

func SubvolumeList(volume string) (result []string) {
	cmd := exec.Command("btrfs", "subvolume", "list", "-a", volume)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		fmt.Errorf("Failed to list subvolumes on root volume: %s", volume)
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

func SubvolumeSnapshot(volume string, sourcesubvolume string, targetsubvolume string, readonly bool) error {
	var cmd *exec.Cmd
	if readonly {
		cmd = exec.Command("btrfs", "subvolume", "snapshot", "-r", path.Join(volume, sourcesubvolume), path.Join(volume, targetsubvolume))
	} else {
		cmd = exec.Command("btrfs", "subvolume", "snapshot", path.Join(volume, sourcesubvolume), path.Join(volume, targetsubvolume))
	}
	return cmd.Run()
}

func SubvolumeDelete(volume string, name string) error {
	cmd := exec.Command("btrfs", "subvolume", "delete", path.Join(volume, name))
	return cmd.Run()
}

func Sync() error {
	cmd := exec.Command("sync")
	return cmd.Run()
}

/* Gets the list of snapshot for a parent subvolume*/
func SubvolumeSnapshotList(volume string, parentsubvolume string) (result []string) {
	cmd := exec.Command("btrfs", "subvolume", "show", path.Join(volume, parentsubvolume))
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		fmt.Errorf("Failed to get subvolume information: %s", path.Join(volume, parentsubvolume))
		return
	} else {
		lines := strings.Split(out.String(), "\n")
		parent_uuid := strings.TrimSpace(strings.Split(lines[2], "uuid:")[1])
		cmd2 := exec.Command("btrfs", "subvolume", "list", "-qus", path.Join(volume, parentsubvolume))
		var out2 bytes.Buffer
		cmd2.Stdout = &out2
		err2 := cmd2.Run()
		if err2 != nil {
			fmt.Errorf("Failed to list snapshots on subvolume : %s", path.Join(volume, parentsubvolume))
		} else {
			lines2 := strings.Split(out2.String(), "\n")
			for _, line2 := range lines2 {
				tokens := strings.Split(line2, " ")
				if len(tokens) == 18 && tokens[13] == parent_uuid {
					result = append(result, tokens[17])
				}
			}
		}
	}
	return
}

func SendReceive(sourcevolume string, sourcesubvolume string, targetvolume string, parentsubvolume string, remotehost string) {
	var cmd1 *exec.Cmd
	var cmd2 *exec.Cmd	
	if(parentsubvolume != "") {
		cmd1 = exec.Command("btrfs","send","-p",path.Join(sourcevolume,parentsubvolume),path.Join(sourcevolume,sourcesubvolume))	
	} else {
		cmd1 = exec.Command("btrfs","send",path.Join(sourcevolume,sourcesubvolume))
	}

	if(remotehost != "") {
		cmd2 = exec.Command("ssh",remotehost,"btrfs","receive",targetvolume)
	} else {
		cmd2 = exec.Command("btrfs","receive",targetvolume)
	}
    	r, w := io.Pipe() 
    	cmd1.Stdout = w
    	cmd2.Stdin = r

    	var b2 bytes.Buffer
    	cmd2.Stdout = &b2

    	cmd1.Start()
    	cmd2.Start()
    	err1 := cmd1.Wait()
    	w.Close()
    	err2 := cmd2.Wait()
    	io.Copy(os.Stdout, &b2)
	if(err1 != nil) {
		fmt.Println("Error Send ",err1)
	}
	if(err2 != nil) {
		fmt.Println("Error Receive ", err2)
	}
}
