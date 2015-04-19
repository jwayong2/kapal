package config

import (
	"fmt"	
	"github.com/hoodiez/kapal/cmds"
	"code.google.com/p/gcfg"
)

type KapalConfig struct {
	Volume map[string]*struct {
		Pool string
		Type string
		Dockerize bool
		Dockername string
		Dockervol string
	}
}

func ReadKapalFile(cfgFile string) KapalConfig {
	var err error
	var cfg KapalConfig
	if cfgFile == "" {
		cfgFile = "./Kapalfile"
	}
	
	err = gcfg.ReadFileInto(&cfg, cfgFile)
	
	if err != nil {
		fmt.Println("Error reading Kapalfile")
	}

	/* This is work in progress, will need to be placed into separate func*/
    	for k := range cfg.Volume {
		fmt.Println("Creating Volumes:", k)
		cmds.CreateVolume(cfg.Volume[k].Pool,k,cfg.Volume[k].Dockerize,cfg.Volume[k].Dockername,cfg.Volume[k].Dockervol)
   	}
	
	return cfg
}

