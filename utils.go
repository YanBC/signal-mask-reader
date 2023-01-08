package main

import (
	"fmt"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

//
// signal num to signal string
//
func getSignalMap() map[int]string {
	ret := make(map[int]string)
	for i := 0; i < 1024; i++ {
		ss := syscall.Signal(i)
		name := unix.SignalName(ss)
		ret[i] = name
	}
	return ret
}

func parseSignalSet(num_mask uint64) []string {
	sigMap := getSignalMap()
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

//
// string processer
//
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

func grapSignalSet(str string, subStr string) (uint64, error) {
	idx := strings.Index(str, subStr)
	if idx == -1 {
		return 0, fmt.Errorf("signal set not found")
	}

	var posStart int
	var posEnd int
	current := idx + len(subStr)
	for {
		char := string(str[current])
		if char == ":" || isWhiteSpace(char) {
			current = current + 1
		} else if isNumeric(char) {
			posStart = current
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
			posEnd = current
			break
		}
	}

	signalSetStr := string(str[posStart:posEnd])
	num, err := strconv.ParseUint(signalSetStr, 16, len(signalSetStr)*4)
	if err != nil {
		return 0, err
	}

	return num, nil
}
