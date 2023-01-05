package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

const (
	VERSION = "1.0.0"
)

func getSignalMap() map[int]string {
	ret := make(map[int]string)
	for i := 0; i < 1024; i++ {
		ss := syscall.Signal(i)
		name := unix.SignalName(ss)
		ret[i] = name
	}
	return ret
}

func isNumeric(str string) bool {
	_, err := strconv.Atoi(str)
	if err != nil {
		if str == "a" || str == "b" || str == "c" || str == "d" || str == "e" || str == "f" {
			return true
		} else {
			return false
		}
	} else {
		return true
	}
}

func isWhiteSpace(char string) bool {
	if char == " " || char == "\n" || char == "\t" {
		return true
	} else {
		return false
	}
}

func grapMask(str string, sub_str string) (uint64, error) {
	idx := strings.Index(str, sub_str)
	if idx == -1 {
		return 0, fmt.Errorf("signal mask not found")
	}

	var start_pos int
	var end_pos int
	current := idx + len(sub_str)
	for {
		char := string(str[current])
		if char == ":" || isWhiteSpace(char) {
			current = current + 1
		} else if isNumeric(char) {
			start_pos = current
			current = current + 1
			break
		} else {
			return 0, fmt.Errorf("%s in string", char)
		}
	}

	for {
		char := string(str[current])
		if isNumeric(char) {
			current = current + 1
		} else {
			end_pos = current
			break
		}
	}

	num_str := string(str[start_pos:end_pos])
	num, err := strconv.ParseUint(num_str, 16, len(num_str)*4)
	if err != nil {
		return 0, err
	}

	return num, nil
}

func parseMask(sigMap map[int]string, num_mask uint64) []string {
	ret := make([]string, 0, 1)

	binary_mask := strconv.FormatUint(num_mask, 2)
	sig := 1
	for idx := len(binary_mask) - 1; idx >= 0; idx-- {
		if string(binary_mask[idx]) == "1" && sigMap[sig] != "" {
			ret = append(ret, sigMap[sig])
		}
		sig++
	}

	return ret
}

func printErrorAndExit(err_msg string) {
	fmt.Println(err_msg)
	os.Exit(-1)
}

func printHelpAndExit() {
	fmt.Println("sigReader Version", VERSION)
	fmt.Printf("Usage: %s -h\n", os.Args[0])
	os.Exit(-2)
}

func printProcessMask(pid *string) {
	status_file_path := filepath.Join("/proc", *pid, "status")
	status_file, err := os.ReadFile(status_file_path)
	if err != nil {
		printErrorAndExit(err.Error())
	}
	status_str := string(status_file)
	// fmt.Println(status_str)

	signal_map := getSignalMap()

	sigPnd := "SigPnd"
	sigBlk := "SigBlk"
	sigIgn := "SigIgn"
	sigCgt := "SigCgt"

	sigPnd_num, err := grapMask(status_str, sigPnd)
	if err != nil {
		printErrorAndExit(err.Error())
	}
	sigPnd_arr := parseMask(signal_map, sigPnd_num)
	fmt.Println("Pending Signal: ", "[", strings.Join(sigPnd_arr, ", "), "]")

	sigBlk_num, err := grapMask(status_str, sigBlk)
	if err != nil {
		printErrorAndExit(err.Error())
	}
	sigBlk_arr := parseMask(signal_map, sigBlk_num)
	fmt.Println("Blocked Signal: ", "[", strings.Join(sigBlk_arr, ", "), "]")

	sigIgn_num, err := grapMask(status_str, sigIgn)
	if err != nil {
		printErrorAndExit(err.Error())
	}
	sigIgn_arr := parseMask(signal_map, sigIgn_num)
	fmt.Println("Ignored Signal: ", "[", strings.Join(sigIgn_arr, ", "), "]")

	sigCgt_num, err := grapMask(status_str, sigCgt)
	if err != nil {
		printErrorAndExit(err.Error())
	}
	sigCgt_arr := parseMask(signal_map, sigCgt_num)
	fmt.Println("Caught Signal: ", "[", strings.Join(sigCgt_arr, ", "), "]")
}

func printMaskParse(mask *string) {
	num, err := strconv.ParseUint(*mask, 16, len(*mask)*4)
	if err != nil {
		printErrorAndExit(err.Error())
	}

	signal_map := getSignalMap()
	sig_arr := parseMask(signal_map, num)
	fmt.Println("Signals: ", "[", strings.Join(sig_arr, ", "), "]")
}

func main() {
	pid := flag.String("pid", "", "process id")
	mask := flag.String("mask", "", "mask string")
	flag.Parse()

	if *pid == "" && *mask == "" {
		printHelpAndExit()
	}

	if *pid != "" {
		printProcessMask(pid)
	} else if *mask != "" {
		printMaskParse(mask)
	}

}
