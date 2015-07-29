package main

import "fmt"

func (z *zpool) NagiosFormat() (message string, exitcode int) {
	if (!z.healthy || z.capacity >= capCritical) && z.faulted > 0 {
		return fmt.Sprintf("CRITICAL: %s %s, capacity: %d%%, faulted: %d\n", z.name, z.status, z.capacity, z.faulted), 2
	} else if (!z.healthy || z.capacity >= capCritical) && z.faulted == 0 {
		return fmt.Sprintf("CRITICAL: %s %s, capacity: %d%%\n", z.name, z.status, z.capacity), 2
	} else if z.capacity >= capWarning {
		return fmt.Sprintf("WARNING: %s %s, capacity: %d%%\n", z.name, z.status, z.capacity), 1
	} else if z.healthy && z.capacity < capWarning && z.faulted == 0 {
		return fmt.Sprintf("OK: %s %s, capacity: %d%%\n", z.name, z.status, z.capacity), 0
	}
	return fmt.Sprintf("UNKNOWN: %s %s, capacity: %d%%, faulted: %d\n", z.name, z.status, z.capacity, z.faulted), 3
}
