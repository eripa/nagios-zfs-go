package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	toolVersion = "0.1.2"
)

var zfsPool string
var capWarning int64
var capCritical int64
var versionCheck bool

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
	output = strings.Trim(output, "\n")
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
	output := runZpoolCommand([]string{"status", z.name})
	err := z.getFaulted(output)
	if err != nil {
		log.Fatal("Error parsing zpool status")
	}
	output = runZpoolCommand([]string{"list", "-H", "-o", "health", z.name})
	err = z.checkHealth(output)
	if err != nil {
		log.Fatal("Error parsing zpool list -H -o health ", z.name)
	}
	output = runZpoolCommand([]string{"list", "-H", "-o", "cap", z.name})
	err = z.getCapacity(output)
	if err != nil {
		log.Fatal("Error parsing zpool capacity")
	}
}

func checkExistance(pool string) (err error) {
	output := runZpoolCommand([]string{"list", pool})
	if strings.Contains(fmt.Sprintf("%s", output), "no such pool") {
		err = errors.New("No such pool")
	}
	return
}

func runZpoolCommand(args []string) string {
	zpoolPath, err := exec.LookPath("zpool")
	if err != nil {
		log.Fatal("Could not find zpool in PATH")
	}
	cmd := exec.Command(zpoolPath, args...)
	out, _ := cmd.CombinedOutput()
	return fmt.Sprintf("%s", out)
}

func main() {
	if versionCheck {
		fmt.Printf("nagios-zfs-go v%s (https://github.com/eripa/nagios-zfs-go)\n", toolVersion)
		os.Exit(0)
	}
	err := checkExistance(zfsPool)
	if err != nil {
		log.Fatal(err)
	}
	z := zpool{name: zfsPool}
	z.getStatus()
	message, exitcode := z.NagiosFormat()
	fmt.Println(message)
	os.Exit(exitcode)
}
