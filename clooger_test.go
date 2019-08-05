package clogger

import (
	"fmt"
	"regexp"
	"testing"
)

func TestConstLevel(t *testing.T) {
	fmt.Println(DebugLevel)
	fmt.Println(WarnLevel)
	fmt.Println(FatalLevel)
}

func TestErrorFileName(t *testing.T) {
	fileName := "xxx.log"
	regex, err := regexp.Compile("(.*)\\.log")
	if err != nil {
		panic(err)
	}

	s := regex.FindAllString(fileName, -1)
	fmt.Println(s)

}

func TestErrorFileName2(t *testing.T) {
	fileName := "xxx.log"
	reg := regexp.MustCompile(`(.*)\.log`)

	prefix := reg.ReplaceAllString(fileName, "$1")
	fmt.Printf("%q\n", prefix)
	fmt.Println(prefix)

	fmt.Println(fileName)
}
