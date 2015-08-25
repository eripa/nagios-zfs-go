package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	toolVersion = "0.2.1"
)

var zfsPool string
var capWarning int64
var capCritical int64
var versionCheck bool
var dumpDirectory string

func init() {
	const (
		defaultPool     = "tank"
		poolUsage       = "what ZFS pool to check"
		defaultWarning  = 70
		warningUsage    = "Capacity warning limit"
		defaultCritical = 80
		criticalUsage   = "Capacity critical limit (80% is considered soft limit of ZFS)"
		versionUsage    = "Display current version"
	)
	flag.StringVar(&zfsPool, "pool", defaultPool, poolUsage)
	flag.StringVar(&zfsPool, "p", defaultPool, poolUsage+" (shorthand)")
	flag.Int64Var(&capWarning, "warning", defaultWarning, warningUsage)
	flag.Int64Var(&capWarning, "w", defaultWarning, warningUsage+" (shorthand)")
	flag.Int64Var(&capCritical, "critical", defaultCritical, criticalUsage)
	flag.Int64Var(&capCritical, "c", defaultCritical, criticalUsage+" (shorthand)")
	flag.StringVar(&dumpDirectory, "statusdir", "/tmp", "Where to look for status dumps (see separate dumper.sh script)")
	flag.BoolVar(&versionCheck, "version", false, versionUsage)
	flag.Parse()
}

type zpool struct {
	name     string
	capacity int64
	healthy  bool
	status   string
	faulted  int64
}

func (z *zpool) checkHealth(output string) (err error) {
	health := strings.Trim(output, "\n")
	if health == "ONLINE" {
		z.healthy = true
	} else if health == "DEGRADED" || health == "FAULTED" {
		z.healthy = false
	} else {
		z.healthy = false // just to make sure
		err = errors.New("Unknown status")
	}
	return err
}

func (z *zpool) getCapacity(output string) (err error) {
	s := strings.Split(output, "%")[0]
	z.capacity, err = strconv.ParseInt(s, 0, 8)
	if err != nil {
		return err
	}
	return err
}

func (z *zpool) getFaulted(output string) (err error) {
	lines := strings.Split(output, "\n")
	z.status = strings.Split(lines[1], " ")[2]
	if z.status == "ONLINE" {
		z.faulted = 0 // assume ONLINE means no faulted/unavailable providers
	} else if z.status == "DEGRADED" || z.status == "FAULTED" {
		var count int64
		for _, line := range lines {
			if (strings.Contains(line, "FAULTED") && !strings.Contains(line, "state:")) || strings.Contains(line, "UNAVAIL") {
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

func (z *zpool) getStatus() {
	basePath := filepath.Join(dumpDirectory, "check_zfs_"+z.name)
	data, err := ioutil.ReadFile(basePath + "_status")
	if err != nil {
		log.Printf("Error: Have you dumped status files to %s with dumper.sh?\n", dumpDirectory)
		log.Fatal(err)
	}
	err = z.getFaulted(fmt.Sprintf("%s", data))
	if err != nil {
		log.Fatal("Error parsing zpool status")
	}
	data, err = ioutil.ReadFile(basePath + "_health")
	if err != nil {
		log.Printf("Error: Have you dumped status (health) files to %s with dumper.sh?\n", dumpDirectory)
		log.Fatal(err)
	}
	err = z.checkHealth(fmt.Sprintf("%s", data))
	if err != nil {
		log.Fatal("Error parsing zpool health")
	}
	data, err = ioutil.ReadFile(basePath + "_capacity")
	if err != nil {
		log.Printf("Error: Have you dumped status (capacity) files to %s with dumper.sh?\n", dumpDirectory)
		log.Fatal(err)
	}
	err = z.getCapacity(fmt.Sprintf("%s", data))
	if err != nil {
		log.Fatal("Error parsing zpool capacity")
	}
}

func main() {
	if versionCheck {
		fmt.Printf("nagios-zfs-go v%s (https://github.com/eripa/nagios-zfs-go)\n", toolVersion)
		os.Exit(0)
	}
	z := zpool{name: zfsPool}
	z.getStatus()
	message, exitcode := z.NagiosFormat()
	fmt.Println(message)
	os.Exit(exitcode)
}
