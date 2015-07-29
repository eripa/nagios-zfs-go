package main

import (
	"testing"
)

func TestCheckHealth(t *testing.T) {
	z := zpool{name: "tank"}

	// Test ONLINE
	err := checkHealth(&z, "ONLINE")

	if err != nil {
		t.Errorf("Error in checkHealth (%s)", err)
	}
	if z.healthy == false {
		t.Errorf("healthy should equal true when given 'ONLINE'")
	}

	// Test FAULED
	err = checkHealth(&z, "FAULTED")
	if err != nil {
		t.Errorf("Error in checkHealth (%s)", err)
	}
	if z.healthy == true {
		t.Errorf("healthy should equal true when given 'FAULTED'")
	}

	// Test DEGRADED
	err = checkHealth(&z, "DEGRADED")
	if err != nil {
		t.Errorf("Error in checkHealth (%s)", err)
	}
	if z.healthy == true {
		t.Errorf("healthy should equal true when given 'DEGRADED'")
	}

	// Test other
	err = checkHealth(&z, "other status")
	if err == nil {
		t.Errorf("other status should throw error in checkHealth (%s)", err)
	}
	if z.healthy == true {
		t.Errorf("healthy should equal false when given unknown input")
	}
}

func TestGetCapacity(t *testing.T) {
	z := zpool{name: "tank"}

	// Test average capacity
	err := getCapacity(&z, "51%")

	if err != nil {
		t.Errorf("Error in getCapacity")
	}
	if z.capacity != 51 {
		t.Errorf("Non-matching integer, should be 51")
	}

	// Test non-integer
	err = getCapacity(&z, "foo")

	if err == nil {
		t.Errorf("Non-integer should produce error in getCapacity")
	}
}

func TestGetFaulted(t *testing.T) {
	z := zpool{
		name:    "tank",
		faulted: -1, // set to -1 to make sure we actually test the all-ONLINE case
	}

	// Test all ONLINE
	err := getFaulted(&z, `  pool: tank
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
	if err != nil {
		t.Errorf("Error in getFaulted")
	}
	if z.faulted != 0 {
		t.Errorf("Incorrect amount of faulted, should be 0.")
	}

	// Test degraded state
	err = getFaulted(&z, `  pool: tank
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
	if err != nil {
		t.Errorf("Error in getFaulted")
	}
	if z.faulted != 2 {
		t.Errorf("Incorrect amount of faulted, should be 2.")
	}

	// Test other output
	err = getFaulted(&z, `  pool: tank
	 state: Oother`)
	if err == nil {
		t.Errorf("Should produce parsing error in getFaulted")
	}
	if z.faulted != 1 {
		t.Errorf("Incorrect amount of faulted, should be 1 when parsing error.")
	}
}
