# nagios-zfs-go

Nagios check for ZFS pool status, written in Go.

Using Go gives the nice benefit of static binaries on different platforms. The only external dependency is 'zpool', which you probably have where you want to use this.

Tested on:

  * SmartOS v0.147+ build: 20150612T210440Z, with 6-disk raidz2 setup
  * Debian Wheezy 7.8 (3.2.0-4-amd64), zfsonlinux 0.6.4.2, with a 8-disk raidz2 setup

## Usage

    $ bin/check_zfs --help
    Usage of bin/check_zfs:
      -c int
            Capacity critical limit (80% is considered soft limit of ZFS) (shorthand) (default 80)
      -critical int
            Capacity critical limit (80% is considered soft limit of ZFS) (default 80)
      -p string
            what ZFS pool to check (shorthand) (default "tank")
      -pool string
            what ZFS pool to check (default "tank")
      -w int
            Capacity warning limit (shorthand) (default 70)
      -warning int
            Capacity warning limit (default 70)

## Example run

In order to separate privileges and avoid root requirement for check_zfs, the check is split into two steps.

### First, dump data using root privileges and the dumper script

    $ bin/check_zfs_dumper.sh

This will create simple text files with status in /tmp (change script if desired). These files will later on be parsed by check_zfs.

This script would presumably be run by crontab every minute or so. Example crontab entry:

    *  *  * *  *  /root/scripts/nagios-zfs-go/dumper.sh

### Next, run the check_zfs command

A typical example for SmartOS "zones" pool:

    $ bin/check_zfs -p zones -w 70 -c 78
    OK: zones ONLINE, capacity: 51%

    $ echo $?
    0

Or with a strict capacity limit:

    $ bin/check_zfs -p zones -w 40 -c 50
    CRITICAL: zones ONLINE, capacity: 51%

    $ echo $?
    2

## Build

I recommend to use Go 1.5 (currently beta2 as of 2015-07-29), to make cross-compilation a lot easier.

SmartOS (x86_64):

    env GOOS=solaris GOARCH=amd64 go build -o bin/check_zfs-solaris

Linux (x86_64):

    env GOOS=linux GOARCH=amd64 go build -o bin/check_zfs-linux

Mac OS X:

    env GOOS=darwin GOARCH=amd64 go build -o bin/check_zfs-mac

## Tests

There are some simple test cases to make sure that no insane results occur. All test cases are based on a raidz2 setup with 6 disks. So perhaps more variants of pool configurations would be good to add..

Run `go test -v` to run the tests with some verbosity.

## bin/zpool

`bin/zpool` is a shell-script that can be used to fake a 'zpool' command on your local development machine where you might not have ZFS installed. It will simply run zpool over SSH on a remote host. Set environment variable ZFSHOST to whatever host you want to remote to.

The script also has some simple sed statments prepared (you will have to remove the hash signs manually) to fake different pool statuses for testing purposes.

## License

The MIT License, see separate LICENSE file for full text.

## Contributing

  * Fork it
  * Create your feature branch (git checkout -b my-new-feature)
  * Commit your changes (git commit -am 'Add some feature')
  * Push to the branch (git push origin my-new-feature)
  * Create new Pull Request
