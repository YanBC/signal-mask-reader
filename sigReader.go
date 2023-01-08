package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	VERSION = "1.0.1"
)

const (
	sigPnd = "SigPnd"
	shdPnd = "ShdPnd"
	sigBlk = "SigBlk"
	sigIgn = "SigIgn"
	sigCgt = "SigCgt"
)

//
// handle error
//
func printErrorAndExit(errMsg string) {
	fmt.Println(errMsg)
	os.Exit(-1)
}

func printHelpAndExit() {
	fmt.Println("sigReader Version", VERSION)
	fmt.Printf("Usage: %s -h\n", os.Args[0])
	os.Exit(-2)
}

//
// main
//
func grapAndParseSignalSet(statusStr string, sigSetName string) {
	signalSet, err := grapSignalSet(statusStr, sigSetName)
	if err != nil {
		fmt.Println(sigSetName, err)
		return
	}
	signalSetArr := parseSignalSet(signalSet)
	fmt.Println(sigSetName, strings.Join(signalSetArr, ", "))
}

func printProcess(pid *string) {
	statusFilePath := filepath.Join("/proc", *pid, "status")
	statusFile, err := os.ReadFile(statusFilePath)
	if err != nil {
		printErrorAndExit(err.Error())
	}
	statusStr := string(statusFile)

	sigSetNames := []string{
		sigPnd, shdPnd, sigBlk, sigIgn, sigCgt,
	}

	for _, sigSetName := range sigSetNames {
		grapAndParseSignalSet(statusStr, sigSetName)
	}
}

func printSignalSet(signalSetStr *string) {
	signalSet, err := strconv.ParseUint(*signalSetStr, 16, len(*signalSetStr)*4)
	if err != nil {
		printErrorAndExit(err.Error())
	}

	signalSetArr := parseSignalSet(signalSet)
	fmt.Println("Signals: ", "[", strings.Join(signalSetArr, ", "), "]")
}

func main() {
	pid := flag.String("pid", "", "process id / thread id")
	signalSetStr := flag.String("parse", "", "parse signal set string")
	flag.Parse()

	if *pid == "" && *signalSetStr == "" {
		printHelpAndExit()
	}

	if *pid != "" {
		printProcess(pid)
	} else if *signalSetStr != "" {
		printSignalSet(signalSetStr)
	}

}
