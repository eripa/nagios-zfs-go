package main

import (
	"strings"
	"testing"
)

func TestNagiosOK(t *testing.T) {
	z := zpool{
		name:    "tank",
		faulted: -1, // set to -1 to make sure we actually test the all-ONLINE case
	}

	// Test all ONLINE
	z.checkHealth("ONLINE")
	z.getCapacity("51%")
	z.getFaulted(`  pool: tank
 state: ONLINE
  scan: scrub repaired 0 in 1h1m with 0 errors on Thu Jan 1 13:37:00 1970
config:

        NAME                       STATE     READ WRITE CKSUM
        zones                      ONLINE       0     0     0
          raidz2-0                 ONLINE       0     0     0
            c0t5000C5006A6E87D9d0  ONLINE       0     0     0
            c0t5000C50024CAAFFCd0  ONLINE       0     0     0
            c0t5000CCA249D27B4Ed0  ONLINE       0     0     0
            c0t5000C5004425F6F6d0  ONLINE       0     0     0
            c0t5000C500652DD0EFd0  ONLINE       0     0     0
            c0t50014EE25A580141d0  ONLINE       0     0     0

errors: No known data errors`)

	message, exitcode := z.NagiosFormat()
	if message != "OK: tank ONLINE, capacity: 51%" {
		t.Errorf("Unexpected Nagios status. Got: %s", message)
	}
	if exitcode != 0 {
		t.Errorf("Unexpected Exit code, got %d, should be 0", exitcode)
	}
}

func TestNagiosWarningCap(t *testing.T) {
	z := zpool{
		name:    "tank",
		faulted: -1, // set to -1 to make sure we actually test the all-ONLINE case
	}
	z.checkHealth("ONLINE")
	z.getCapacity("78%")
	z.getFaulted(`  pool: tank
	 state: ONLINE
	  scan: scrub repaired 0 in 1h1m with 0 errors on Thu Jan 1 13:37:00 1970
	config:

	        NAME                       STATE     READ WRITE CKSUM
	        zones                      ONLINE       0     0     0
	          raidz2-0                 ONLINE       0     0     0
	            c0t5000C5006A6E87D9d0  ONLINE       0     0     0
	            c0t5000C50024CAAFFCd0  ONLINE       0     0     0
	            c0t5000CCA249D27B4Ed0  ONLINE       0     0     0
	            c0t5000C5004425F6F6d0  ONLINE       0     0     0
	            c0t5000C500652DD0EFd0  ONLINE       0     0     0
	            c0t50014EE25A580141d0  ONLINE       0     0     0

	errors: No known data errors`)
	message, exitcode := z.NagiosFormat()
	if message != "WARNING: tank ONLINE, capacity: 78%" {
		t.Errorf("Unexpected Nagios status. Got: %s", message)
	}
	if exitcode != 1 {
		t.Errorf("Unexpected Exit code, got %d, should be 1", exitcode)
	}
}

func TestNagiosCriticalCap(t *testing.T) {

	z := zpool{
		name:    "tank",
		faulted: -1, // set to -1 to make sure we actually test the all-ONLINE case
	}
	z.checkHealth("ONLINE")
	z.getCapacity("88%")
	z.getFaulted(`  pool: tank
	 state: ONLINE
	  scan: scrub repaired 0 in 1h1m with 0 errors on Thu Jan 1 13:37:00 1970
	config:

	        NAME                       STATE     READ WRITE CKSUM
	        zones                      ONLINE       0     0     0
	          raidz2-0                 ONLINE       0     0     0
	            c0t5000C5006A6E87D9d0  ONLINE       0     0     0
	            c0t5000C50024CAAFFCd0  ONLINE       0     0     0
	            c0t5000CCA249D27B4Ed0  ONLINE       0     0     0
	            c0t5000C5004425F6F6d0  ONLINE       0     0     0
	            c0t5000C500652DD0EFd0  ONLINE       0     0     0
	            c0t50014EE25A580141d0  ONLINE       0     0     0

	errors: No known data errors`)
	message, exitcode := z.NagiosFormat()
	if message != "CRITICAL: tank ONLINE, capacity: 88%" {
		t.Errorf("Unexpected Nagios status. Got: %s", message)
	}
	if exitcode != 2 {
		t.Errorf("Unexpected Exit code, got %d, should be 2", exitcode)
	}

}

func TestNagiosCriticalDegraded(t *testing.T) {
	z := zpool{
		name:    "tank",
		faulted: -1, // set to -1 to make sure we actually test the all-ONLINE case
	}
	z.checkHealth("DEGRADED")
	z.getCapacity("32%")
	z.getFaulted(`  pool: tank
	 state: DEGRADED
	  scan: scrub repaired 0 in 1h1m with 0 errors on Thu Jan 1 13:37:00 1970
	config:

	        NAME                       STATE     READ WRITE CKSUM
	        zones                      DEGRADED     0     0     0
	          raidz2-0                 ONLINE       0     0     0
	            c0t5000C5006A6E87D9d0  FAULTED      0     0     0
	            c0t5000C50024CAAFFCd0  ONLINE       0     0     0
	            c0t5000CCA249D27B4Ed0  ONLINE       0     0     0
	            c0t5000C5004425F6F6d0  UNAVAIL      0     0     0
	            c0t5000C500652DD0EFd0  ONLINE       0     0     0
	            c0t50014EE25A580141d0  ONLINE       0     0     0

	errors: No known data errors`)
	message, exitcode := z.NagiosFormat()
	if message != "CRITICAL: tank DEGRADED, capacity: 32%, faulted: 2" {
		t.Errorf("Unexpected Nagios status. Got: %s", message)
	}
	if exitcode != 2 {
		t.Errorf("Unexpected Exit code, got %d, should be 2", exitcode)
	}
}

func TestNagiosCriticalFaulted(t *testing.T) {
	z := zpool{
		name:    "tank",
		faulted: -1, // set to -1 to make sure we actually test the all-ONLINE case
	}
	z.checkHealth("FAULTED")
	z.getCapacity("43%")
	z.getFaulted(`  pool: tank
	 state: FAULTED
	  scan: scrub repaired 0 in 1h1m with 0 errors on Thu Jan 1 13:37:00 1970
	config:

	        NAME                       STATE     READ WRITE CKSUM
	        zones                      DEGRADED     0     0     0
	          raidz2-0                 ONLINE       0     0     0
	            c0t5000C5006A6E87D9d0  FAULTED      0     0     0
	            c0t5000C50024CAAFFCd0  ONLINE       0     0     0
	            c0t5000CCA249D27B4Ed0  FAULTED      0     0     0
	            c0t5000C5004425F6F6d0  FAULTED      0     0     0
	            c0t5000C500652DD0EFd0  ONLINE       0     0     0
	            c0t50014EE25A580141d0  ONLINE       0     0     0

	errors: No known data errors`)
	message, exitcode := z.NagiosFormat()
	if message != "CRITICAL: tank FAULTED, capacity: 43%, faulted: 3" {
		t.Errorf("Unexpected Nagios status. Got: %s", message)
	}
	if exitcode != 2 {
		t.Errorf("Unexpected Exit code, got %d, should be 2", exitcode)
	}
}

func TestNagiosOutput(t *testing.T) {
	nagiosStatus := map[string]int{
		"OK:":       0,
		"WARNING:":  1,
		"CRITICAL:": 2,
		"UNKNOWN:":  3,
	}
	pools := map[string]zpool{
		"OK: tank ONLINE, capacity: 43%": {
			name:     "tank",
			faulted:  0,
			capacity: 43,
			healthy:  true,
			status:   "ONLINE",
		},
		"WARNING: zones ONLINE, capacity: 78%": {
			name:     "zones",
			faulted:  0,
			capacity: 78,
			healthy:  true,
			status:   "ONLINE",
		},
		"CRITICAL: tank ONLINE, capacity: 83%": {
			name:     "tank",
			faulted:  0,
			capacity: 83,
			healthy:  true,
			status:   "ONLINE",
		},
		"CRITICAL: tank DEGRADED, capacity: 43%, faulted: 1": {
			name:     "tank",
			faulted:  1,
			capacity: 43,
			healthy:  false,
			status:   "DEGRADED",
		},
		"CRITICAL: tank FAULTED, capacity: 13%, faulted: 2": {
			name:     "tank",
			faulted:  2,
			capacity: 13,
			healthy:  false,
			status:   "FAULTED",
		},
		"UNKNOWN: zones OTHERSTATE, capacity: -1%, faulted: 1": {
			name:     "zones",
			faulted:  1,
			capacity: -1,
			healthy:  true,
			status:   "OTHERSTATE",
		},
	}

	for status, pool := range pools {
		message, exitcode := pool.NagiosFormat()
		if message != status {
			t.Errorf("Unexpected Nagios status. Got: '%s', should be: '%s'", message, status)
		}
		s := strings.Fields(message)[0]
		if exitcode != nagiosStatus[s] {
			t.Errorf("Unexpected Exit code, got %d, should be %d", exitcode, nagiosStatus[s])
		}
	}
}
