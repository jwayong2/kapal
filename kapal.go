package main

import (
	"fmt"
	"github.com/hoodiez/kapal/btrfs"
)

func main() {
	volumes := btrfscmd.SubvolumeList("/var/lib/docker")
	for _,vol := range volumes {
		fmt.Println(vol)
	}
}
