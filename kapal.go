package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"bytes"
	"strings"
	"github.com/codegangsta/cli"
	"github.com/hoodiez/kapal/btrfs"
	"github.com/hoodiez/kapal/sysfs"
)

func CreateVolume (pool string, name string, dockerize bool, dockername string, dockervol string) {
	err := btrfscmd.SubvolumeCreate(pool,name)
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
	err := btrfscmd.SubvolumeSnapshot(pool,source,target,readonly)
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

func main() {
	version := "0.0.1"
	app := cli.NewApp()
	app.Name = "kapal"
	app.Usage = "Linux Container data orchestration tool"
	app.Version = version 
	app.Authors = []cli.Author{cli.Author{Name: "hoodiez", Email:"https://github.com/hoodiez"},cli.Author{Name: "jteso",Email: "https://github.com/jteso"}}
	app.Commands = []cli.Command {
		{
			Name: "devices",
			Usage: "Exposes information about your system device tree"
			Subcommands: []cli.Command {
				Name: "list",
				Usage: "List of all attached storage devices",
				Action: func(c *cli.Context){
						sysfs.ListDevicesCmd()
				}
			}
		},
		{
			Name: "pool",
			Aliases: []string{"p"},
			Usage: "Manage storage pools",
			Subcommands: []cli.Command {
			 {
			 	Name: "create",
				Usage: "Create new storage pool from a device",
				Action: func(c *cli.Context) {
					fmt.Println("Add Storage Pool: ", c.String("device"), c.String("mount"))
				},
				Flags: []cli.Flag {
				 cli.StringFlag{
				  Name: "device, d",
				  Usage: "device name, e.g. /dev/sdb",
				 },
				 cli.StringFlag{
				  Name: "mount, m",
				  Usage: "mount point path of the storage pool, e.g. /var/lib/kapal",
				 },
				},
			 },
			 {
				Name: "list",
				Usage: "List storage pools that can be used or managed by kapal",
				Action: func(c *cli.Context) {
                                        fmt.Println("List Storage Pool: ")
                                },
			 },
			},
		},
		{
			Name: "volume",
			Aliases: []string{"vol"},
			Usage: "Manage volumes in a storage pool, such as create, clone, remove, etc.",
			Subcommands: []cli.Command {
			{
                                Name: "create",
                                Usage: "Create a new volume in a pool",
                                Action: func(c *cli.Context) {
                                        fmt.Println("Create Volume: ", c.String("pool"), c.String("name"))
					CreateVolume(c.String("pool"),c.String("name"),c.Bool("dockerize"),c.String("dockername"),c.String("dockervol"))
                                },
                                Flags: []cli.Flag {
                                 cli.StringFlag{
                                  Name: "pool, p",
                                  Usage: "pool mount point, e.g. /var/lib/kapal",
                                 },
                                 cli.StringFlag{
                                  Name: "name, n",
                                  Usage: "name of the volume, e.g. vol01",
                                 },
				 cli.BoolFlag{
                                  Name: "dockerize, d",
                                  Usage: "Also create a docker data volume container, default is false",
                                 },
                                 cli.StringFlag{
                                  Name: "dockername",
                                  Usage: "Name of docker data volume container, default will use docker automatic naming",
                                 },
                                 cli.StringFlag{
                                  Name: "dockervol",
                                  Usage: "Docker Volume path mounted inside the container, default is /data",
                                 },
                                },
                         },
                         {
                                Name: "list",
                                Usage: "List volumes in a pool",
                                Action: func(c *cli.Context) {
                                        fmt.Println("List Volumes: ", c.String("pool"))
                                },
				Flags: []cli.Flag {
                                 cli.StringFlag{
                                  Name: "pool, p",
                                  Usage: "pool mount point, e.g. /var/lib/kapal",
                                 },
                                },
                         },
		         {
                                Name: "remove",
                                Usage: "Remove a volume in a pool",
                                Action: func(c *cli.Context) {
                                        fmt.Println("Remove Volume: ", c.String("pool"))
                                },      
                                Flags: []cli.Flag {
                                 cli.StringFlag{
                                  Name: "pool, p",
                                  Usage: "pool mount point, e.g. /var/lib/kapal",
                                 },
                                 cli.StringFlag{
                                  Name: "name, n",
                                  Usage: "name of the volume, e.g. vol01",
                                 },
                                },
                         },
			 {
				Name: "clone",
				Usage: "Clone a volume into a another volume in a pool",
				Action: func(c *cli.Context) {
                                        fmt.Println("Clone Volume: ", c.String("pool"),c.String("source"),c.String("target"))
					CloneVolume(c.String("pool"),c.String("source"),c.String("target"),c.Bool("readonly"),c.Bool("dockerize"),c.String("dockername"),c.String("dockervol"))
                                },
                                Flags: []cli.Flag {
                                 cli.StringFlag{
                                  Name: "pool, p",
                                  Usage: "pool mount point, e.g. /var/lib/kapal",
                                 },
                                 cli.StringFlag{
                                  Name: "source, s",
                                  Usage: "source volume name",
                                 },
				 cli.StringFlag{
				  Name: "target, t",
				  Usage: "target clone volume name",
				 },
				 cli.BoolFlag{
				  Name: "readonly, r",
				  Usage: "make target clone volume readonly, default is false",
                                 },
                                 cli.BoolFlag{
                                  Name: "dockerize, d",
                                  Usage: "Also create a docker data volume container, default is false",
                                 },
                                 cli.StringFlag{
                                  Name: "dockername",
                                  Usage: "Name of docker data volume container, default will use docker automatic naming",
                                 },
                                 cli.StringFlag{
                                  Name: "dockervol",
                                  Usage: "Docker Volume path mounted inside the container, default is /data",
                                 },
				},
			 },
			 {
				Name: "backup",
				Usage: "Backup a volume from one pool to another local or remote pool",
				Action: func(c *cli.Context) {
                                        fmt.Println("Clone Volume: ", c.String("pool"),c.String("source"),c.String("target"))
                                },
				Flags: []cli.Flag {
                                 cli.StringFlag{
                                  Name: "pool, p",
                                  Usage: "pool mount point, e.g. /var/lib/kapal",
                                 },
                                 cli.StringFlag{
                                  Name: "targetpool, t",
                                  Usage: "target pool name",
                                 },
				 cli.StringFlag{
                                  Name: "n, name",
                                  Usage: "volume name",
                                 },
				 cli.StringFlag{
                                  Name: "remote, r",
                                  Usage: "remote host or ip address",
                                 },
				 cli.StringFlag{
				  Name: "clonename, c",
				  Usage: "specific backup clone name",
				 },	
				},	
			 },
			},		
		},
	}

	app.Run(os.Args)

	/*Testing Btrfs*/
}
