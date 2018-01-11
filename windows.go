package main

import (
    "strings"
    "strconv"
    "log"
)

func windows_format_name(location string, prev string) (name string) {
    // so apparently you need the space to be ascii code... rip 90 minutes
    // get everything after prev, in that result, replace all spaces w/ ascii code
    var spaces string
    if prev == "<DIR>" {
        spaces = "          "
    } else {
        spaces = " "
    }
    name = strings.Replace(strings.Split(location, (prev + spaces))[1], " ", "\x20", -1)
    return name
}

// Helper that execs ls -l to generate an Info struct of everything inside the given directory
// for files it just saves the info
// for directories it recurses to get the sum size of all items inside that directory
func windows_scan_dir_contents(location string) ([]*Info, uint64) {
    var total_size uint64 = 0

    /*
    // command to get rid of any alias to ensure format
    cmd := exec.Command(dir, flags, location)
    out, err := cmd.Output()
    // gracefully continue - return errored info struct list
    if err != nil {
        bad_dirs = append(bad_dirs, location)

        return []*Info{
            &Info{
                directory : true,
                owner     : "ERR",
                group     : "ERR",
                size      : 0,
                name      : "ERROR - NO ACCESS TO PARENT",
                children  : []*Info{},
                parent    : nil,
            },
        }, 0
    }

    */

    out := []byte(` Volume in drive C has no label.
 Volume Serial Number is F4AC-9851

 Directory of C:\

09/02/2015  12:41 PM    <DIR>          $SysReset
05/30/2016  06:22 PM                93 HaxLogs.txt
05/07/2016  02:58 AM    <DIR>          PerfLogs
05/22/2016  07:55 PM    <DIR>          Program Files
05/31/2016  11:30 AM    <DIR>          Program Files (x86)
07/30/2015  04:32 PM    <DIR>          Temp
05/22/2016  07:55 PM    <DIR>          Users
05/22/2016  08:00 PM    <DIR>          Windows
05/22/2016  09:50 PM    <DIR>          Windows.old
               1 File(s)             93 bytes
               8 Dir(s)  18,370,433,024 bytes free`)
    if len(out) == 0 {
        // return an empty slice if nothing in folder
        return []*Info{}, 0
    }

    // split and ignore first five and last two info lines the two newlines at the end
    spl := strings.Split(string(out), "\n")
    take := spl[5 : (len(spl) - 4)]

    var contents []*Info

    for _, line := range take {
        fields := strings.Fields(line)

        // index 3 is right before the file name
        name := windows_format_name(line, fields[3])

        log.Printf("%#v", name)


        directory := (fields[3] == "<DIR>")

        var size uint64
        var err error
        var children []*Info

        if name == "." || name == ".." {
            continue
        } else if directory {
            children, size = []*Info{}, 69 // windows_scan_dir_contents(location + "\\" + name)
        } else {
            size, err = strconv.ParseUint(fields[3], 10, 64)
            if err != nil {
                log.Fatal(err)
            }
            children = nil // explicit
        }
        // add to total size for current directory
        total_size += size

        info := Info{   
            directory : directory,
            owner     : "sad",
            group     : "boiz",
            size      : size,
            name      : name,
            children  : children,
            parent    : nil,
        }

        // set parent for children - not sure if this is the best way
        // use index cause range copies the values, when we need to modify
        for i, _ := range info.children {
            info.children[i].parent = &info
        }

        contents = append(contents, &info)
    }

    // pprint_children(contents)
    // fmt.Printf("Size of %s: %d\n==========================\n", location, total_size)

    // fmt.Printf("CONTENTS:%#v\n\n", contents)
    return contents, total_size
}

// Scan a directory - takes path to directory to scan
// Produces an Info struct with calculated size and directory attributes
func windows_scan_dir(location string) Info {
    // replace spaces with ascii code - to make it work with exec
    name := strings.Replace(location, " ", "\x20", -1)
    children, size := windows_scan_dir_contents(name)
    
    info := Info{   
        directory : true,
        owner     : "n/a",
        group     : "n/a",
        size      : size,
        name      : name,
        children  : children,
        parent    : nil,
    }

    for i, _ := range info.children {
        info.children[i].parent = &info
    }

    return info
}