package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ConfInit() {
	confFilePath, err := os.Getwd()
	if err != nil {
		fmt.Printf("get pwd error:%v", err.Error())
		os.Exit(Exit_CmdLineParaErr)
	}
	if confFilePath == "" {
		fmt.Printf("\x1b[%dm[info]\x1b[0m Use '-h' to view help information\n", 43)
		os.Exit(Exit_CmdLineParaErr)
	}
	_, err = os.Stat(confFilePath)
	if err != nil {
		fmt.Printf("\x1b[%dm[err]\x1b[0m The '%v' file does not exist\n", 31, confFilePath)
		os.Exit(Exit_ConfFileNotExist)
	}
	path := filepath.Join(confFilePath, "conf.yaml")
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("\x1b[%dm[err]\x1b[0m The '%v' file read error\n", 41, confFilePath)
		os.Exit(Exit_ConfFileTypeError)
	}
	err = yaml.Unmarshal(yamlFile, Data)
	if err != nil {
		fmt.Printf("\x1b[%dm[err]\x1b[0m The '%v' file format error\n", 41, confFilePath)
		os.Exit(Exit_ConfFileFormatError)
	}
}
