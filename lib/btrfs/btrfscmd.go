package btrfscmd

import (
	"os/exec"
	"path"
)

func SubvolumeCreate (volume string, name string) error {
	cmd := exec.Command("btrfs", "subvolume", "create", path.Join(volume, name))
	return cmd.Run()
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


