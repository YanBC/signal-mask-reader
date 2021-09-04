package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
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
		if string(binary_mask[idx]) == "1" {
			ret = append(ret, sigMap[sig])
		}
		sig++
	}

	return ret
}

func main() {
	pid := flag.String("pid", "", "process id")
	flag.Parse()

	if *pid == "" {
		log.Fatal("You must specify a process id")
	}

	status_file_path := filepath.Join("/proc", *pid, "status")
	status_file, err := os.ReadFile(status_file_path)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err.Error())
	}
	sigPnd_arr := parseMask(signal_map, sigPnd_num)
	fmt.Println("Pending Signal: ", sigPnd_arr)

	sigBlk_num, err := grapMask(status_str, sigBlk)
	if err != nil {
		log.Fatal(err.Error())
	}
	sigBlk_arr := parseMask(signal_map, sigBlk_num)
	fmt.Println("Blocked Signal: ", sigBlk_arr)

	sigIgn_num, err := grapMask(status_str, sigIgn)
	if err != nil {
		log.Fatal(err.Error())
	}
	sigIgn_arr := parseMask(signal_map, sigIgn_num)
	fmt.Println("Ignored Signal: ", sigIgn_arr)

	sigCgt_num, err := grapMask(status_str, sigCgt)
	if err != nil {
		log.Fatal(err.Error())
	}
	sigCgt_arr := parseMask(signal_map, sigCgt_num)
	fmt.Println("Caught Signal: ", sigCgt_arr)
}
