package main

import "testing"

func TestNagiosOK(t *testing.T) {
	z := zpool{
		name:    "tank",
		faulted: -1, // set to -1 to make sure we actually test the all-ONLINE case
	}

	// Test all ONLINE
	checkHealth(&z, "ONLINE")
	getCapacity(&z, "51%")
	getFaulted(&z, `  pool: tank
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
	checkHealth(&z, "ONLINE")
	getCapacity(&z, "78%")
	getFaulted(&z, `  pool: tank
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
	checkHealth(&z, "ONLINE")
	getCapacity(&z, "88%")
	getFaulted(&z, `  pool: tank
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
	checkHealth(&z, "DEGRADED")
	getCapacity(&z, "32%")
	getFaulted(&z, `  pool: tank
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
	checkHealth(&z, "FAULTED")
	getCapacity(&z, "43%")
	getFaulted(&z, `  pool: tank
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
