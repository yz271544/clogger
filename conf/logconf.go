package conf

import (
	"github.com/pelletier/go-toml"
)

type cloggerConfig struct {
	FileConf    FileConf
	ConsoleConf ConsoleConf
}

func NewCloggerConfig() *cloggerConfig {
	c := cloggerConfig{}
	if err := c.ParserConfig(); err != nil {
		panic(err)
	}
	return &c
}

type FileConf struct {
	FilePath string `toml:"file_path"`
	FileName string `toml:"file_name"`
	LogLevel string `toml:"log_level"`
}

type ConsoleConf struct {
	LogLevel string `toml:"log_level"`
}

func (c *cloggerConfig) ParserConfig() error {
	tree, err := toml.LoadFile("../conf/application.toml")
	if err != nil {
		return err
	}

	var fileConf FileConf
	var consoleConf ConsoleConf

	fileTree := tree.Get("clogger.file")
	err = fileTree.(*toml.Tree).Unmarshal(&fileConf)
	if err != nil {
		return err
	}

	consoleTree := tree.Get("clogger.console")
	err = consoleTree.(*toml.Tree).Unmarshal(&consoleConf)
	if err != nil {
		return err
	}

	c.FileConf = fileConf
	c.ConsoleConf = consoleConf
	return nil
}
