package cmds

import (
	"bytes"
	"fmt"
	"github.com/hoodiez/kapal/btrfs"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func CreateVolume(pool string, name string, dockerize bool, dockername string, dockervol string) {
	err := btrfs.SubvolumeCreate(pool, name)
	if err != nil {
		fmt.Println("Failed creating Volume in the file system")
	} else {
		/*Todo: should be using Docker API and or it's own wrapper Lib*/
		if dockerize {
			var cmd *exec.Cmd
			var containervol string
			var out bytes.Buffer
			if dockervol != "" {
				containervol = dockervol
				if strings.HasPrefix(containervol, "/") == false {
					containervol = "/" + containervol
				}
			} else {
				containervol = "/data"
			}
			if dockername != "" {
				cmd = exec.Command("docker", "create","-l","kapal=true", "-v", path.Join(pool, name)+":"+containervol, "--name", dockername, "ubuntu")
			} else {
				cmd = exec.Command("docker", "create", "-v", path.Join(pool, name)+":"+containervol, "ubuntu")
			}
			cmd.Stdout = &out
			err2 := cmd.Run()

			if err2 != nil {
				fmt.Println("Error creating Docker Data Volume Container")
			} else {
				fmt.Print(out.String())
			}
		}
	}
}

func CloneVolume(pool string, source string, target string, readonly bool, dockerize bool, dockername string, dockervol string) {
	err := btrfs.SubvolumeSnapshot(pool, source, target, readonly)
	if err != nil {
		fmt.Println("Failed cloning Volume in the file system")
	} else {
		/*Todo: should be using Docker API and or it's own wrapper Lib*/
		if dockerize {
			var cmd *exec.Cmd
			var containervol string
			var out bytes.Buffer

			if dockervol != "" {
				containervol = dockervol
				if strings.HasPrefix(containervol, "/") == false {
					containervol = "/" + containervol
				}
			} else {
				containervol = "/data"
			}
			if dockername != "" {
				cmd = exec.Command("docker", "create","-l","kapal=true", "-v", path.Join(pool, target)+":"+containervol, "--name", dockername, "ubuntu")
			} else {
				cmd = exec.Command("docker", "create", "-v", path.Join(pool, target)+":"+containervol, "ubuntu")
			}
			cmd.Stdout = &out
			err2 := cmd.Run()

			if err2 != nil {
				fmt.Println("Error creating Docker Data Volume Container")
			} else {
				fmt.Print(out.String())
			}

		}
	}
}

func BackupVolume(sourcepool string, targetpool string, volume string, remotehost string) {
	timestamp := time.Now().Local().Format("20060102150405")
	backupname := strings.Join([]string{volume,timestamp},"_")
	parentname := strings.Join([]string{volume,"0"},"_")
	if _, err := os.Stat(path.Join(sourcepool,parentname)); err == nil {
		CloneVolume(sourcepool,volume,backupname,true,false,"","")
		btrfs.Sync()
		btrfs.SendReceive(sourcepool, backupname, targetpool, parentname, remotehost)
	} else {
		CloneVolume(sourcepool,volume,parentname, true, false, "", "")
		btrfs.Sync()
		btrfs.SendReceive(sourcepool, parentname, targetpool, "", remotehost)
	}
}
