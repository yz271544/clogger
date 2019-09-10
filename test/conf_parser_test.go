package main

import (
	"fmt"
	"github.com/kataras/golog"
	"github.com/pelletier/go-toml"
	"github.com/yz271544/clogger/conf"
	"strings"
	"testing"
)

func TestParserConfigMoni(t *testing.T) {
	tree, err := toml.LoadFile("../conf/application.toml")
	if err != nil {
		panic(err)
	}

	var fileConf conf.FileConf
	var consoleConf conf.ConsoleConf

	fileTree := tree.Get("clogger.file")
	err = fileTree.(*toml.Tree).Unmarshal(&fileConf)
	if err != nil {
		panic(err)
	}
	fmt.Println("=====================================")
	fmt.Println("fileConf.FileName:", fileConf.FileName)
	fmt.Println("fileConf.FilePath:", fileConf.FilePath)
	fmt.Println("fileConf.LogLevel:", fileConf.LogLevel)

	consoleTree := tree.Get("clogger.console")
	err = consoleTree.(*toml.Tree).Unmarshal(&consoleConf)
	if err != nil {
		panic(err)
	}
	fmt.Println("=====================================")
	fmt.Println("consoleConf.LogLevel:", consoleConf.LogLevel)

}

func TestParserConfigSelf(t *testing.T) {
	config := conf.NewCloggerConfig()

	fmt.Println(config.ConsoleConf)
	fmt.Println(config.FileConf)
}

func TestParserGologLevel(t *testing.T) {
	level := golog.ParseLevel(strings.ToLower("debug"))
	fmt.Println(level)
	level = golog.ParseLevel(strings.ToLower("DEBUG"))
	fmt.Println(level)
	level = golog.ParseLevel(strings.ToLower("info"))
	fmt.Println(level)
	level = golog.ParseLevel(strings.ToLower("INFO"))
	fmt.Println(level)
	level = golog.ParseLevel(strings.ToLower("error"))
	fmt.Println(level)
	level = golog.ParseLevel(strings.ToLower("ERROR"))
	fmt.Println(level)
	level = golog.ParseLevel(strings.ToLower("WARN"))
	fmt.Println(level)
	level = golog.ParseLevel(strings.ToLower("warn"))
	fmt.Println(level)
	level = golog.ParseLevel(strings.ToLower("FATAL"))
	fmt.Println(level)
	level = golog.ParseLevel(strings.ToLower("fatal"))
	fmt.Println(level)
}
