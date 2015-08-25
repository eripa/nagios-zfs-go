#!/bin/sh
#
# This script should be run by root in Crontab. It will simply run the zpool commands
# and dump them to temporary files. These files will the be read by nagios-zfs-go
#
# This will enable nagios-zfs-go to be run without root privileges, which isn't really
# a sound thing to do for a network service.
#
# A very simple example cron job:
#
# Minute  Hour   Day of Month  Month             Day of Week      Command
# (0-59)  (0-23) (1-31)        (1-12 or Jan-Dec) (0-6 or Sun-Sat)
# *       *      *             *                 *                /root/scripts/nagios-zfs-go/dumper.sh
#

DUMP_LOCATION="/tmp"

# Check for zpool first
command -v zpool > /dev/null 2>&1 || { echo >&2 "Require zpool, but it's not available in PATH. Aborting."; exit 1; }

writeStatusFile() {
    POOL="$1"
    FILENAME="$2"
    COMMAND="$3"
    zpool $COMMAND "$POOL" > "${DUMP_LOCATION%/}/check_zfs_${POOL}_${FILENAME}.$$" 2> /dev/null
    if [ "$?" -eq 0 ] ; then
        # Atomically update the status file
        mv "${DUMP_LOCATION%/}/check_zfs_${POOL}_${FILENAME}.$$" "${DUMP_LOCATION%/}/check_zfs_${POOL}_${FILENAME}"
    else
        echo >&2 "Error running $COMMAND"
        rm -f "${DUMP_LOCATION%/}/check_zfs_${POOL}_${FILENAME}.$$" # Clean up
        exit 1
    fi
}

if [ "$#" -eq "1" ] ; then
    # A single pool provided, use only that one.
    ZPOOLS="$1"
    zpool list -H "$ZPOOLS" > /dev/null 2>&1 || { echo >&2 "Given zpool name doesn't exist. Aborting."; exit 1; }
else
    # Use all pools
    ZPOOLS=`zpool list -H | awk '{print $1}'`
fi

for pool in "$ZPOOLS" ; do
    writeStatusFile "$pool" "status" "status"
    writeStatusFile "$pool" "capacity" "list -H -o cap"
    writeStatusFile "$pool" "health" "list -H -o health"
done

