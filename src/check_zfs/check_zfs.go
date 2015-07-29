package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	VERSION = "0.0.1"
)

// zpool list -H -o name,cap %s
// zpool list -H -o health %s

type zpool struct {
	name     string
	capacity int64
	healthy  bool
	faulted  int64
}

func checkHealth(z *zpool, output string) (err error) {
	if output == "ONLINE" {
		z.healthy = true
	} else if output == "DEGRADED" || output == "FAULTED" {
		z.healthy = false
	} else {
		z.healthy = false // just to make sure
		err = errors.New("Unknown status")
	}
	return err
}

func getCapacity(z *zpool, output string) (err error) {
	s := strings.Split(output, "%")[0]
	z.capacity, err = strconv.ParseInt(s, 0, 8)
	if err != nil {
		return err
	}
	return err
}

func getFaulted(z *zpool, output string) (err error) {
	lines := strings.Split(output, "\n")
	status := strings.Split(lines[1], " ")[2]
	if status == "ONLINE" {
		z.faulted = 0 // assume ONLINE means no faulted/unavailable providers
	} else if status == "DEGRADED" {
		var count int64
		for _, line := range lines {
			if strings.Contains(line, "FAULTED") || strings.Contains(line, "UNAVAIL") {
				count = count + 1
			}
		}
		z.faulted = count
	} else {
		z.faulted = 1 // fake faulted if there is a parsing error
		err = errors.New("Error parsing faulted/unavailable disks")
	}
	return
}

func runZpoolCommand() {

}

func main() {
	fmt.Printf("Hello from check_zfs v%s\n", VERSION)
}
