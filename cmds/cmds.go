package cmds

import (
        "fmt"
        "os/exec"
        "path"
        "bytes"
        "strings"
        "github.com/hoodiez/kapal/btrfs"
)

func CreateVolume (pool string, name string, dockerize bool, dockername string, dockervol string) {
        err := btrfs.SubvolumeCreate(pool,name)
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
                                        containervol = "/"+containervol
                                }
                        } else {
                                containervol = "/data"
                        }
                        if dockername != "" {
                                cmd = exec.Command("docker","create","-v",path.Join(pool,name)+":"+containervol,"--name",dockername,"ubuntu")
                        } else {
                                cmd = exec.Command("docker","create","-v",path.Join(pool,name)+":"+containervol,"ubuntu")
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

func CloneVolume (pool string, source string, target string, readonly bool, dockerize bool, dockername string, dockervol string) {
        err := btrfs.SubvolumeSnapshot(pool,source,target,readonly)
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
                                        containervol = "/"+containervol
                                }
                        } else {
                                containervol = "/data"
                        }
                        if dockername != "" {
                                cmd = exec.Command("docker","create","-v",path.Join(pool,target)+":"+containervol,"--name",dockername,"ubuntu")
                        } else {
                                cmd = exec.Command("docker","create","-v",path.Join(pool,target)+":"+containervol,"ubuntu")
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
