package main

import (
    "os"
    "runtime"
    "fmt"
    "log"
    "flag"
)

// Holds directories that produce an error when accessing
var bad_dirs []string

// flags for the ls/dir command
var flags string

// Print given number of tabs before given printf format string and argument
// Pretty hacky... not very robust - only takes 1 arg...
// but pretty printing will prolly go away in the future anyway
func tab_print(format string, tabs int, arg interface{}) {
    for i := 0; i < tabs; i++ {
        fmt.Printf("\t")
    }

    if (arg == nil) {
        fmt.Printf(format)
    } else {
        fmt.Printf(format, arg)
    }
}

// Pretty print, with tabs and newlines, the contents of an Info struct, recursing on children
// its actually not that pretty - esp for deep recurses
func pprint_children(contents []*Info, tabs int) {
    tab_print("[\n", tabs, nil)
    for _, c := range contents {
        tab_print("\t{\n", tabs, nil)
        tab_print("\tdirectory : %t,\n", tabs, c.directory)
        tab_print("\towner : %s,\n", tabs, c.owner)
        tab_print("\tgroup : %s,\n", tabs, c.group)
        tab_print("\tsize : %d,\n", tabs, c.size)
        tab_print("\tname : %s,\n", tabs, c.name)

        if c.parent == nil {
            tab_print("\tparent (deref'ed) : %s,\n", tabs, "None")            
        } else {
            tab_print("\tparent (deref'ed) : %s,\n", tabs, (*c.parent).name)
        }

        if len(c.children) == 0 {
            tab_print("\tchildren : []\n", tabs, nil)
        } else {
            tab_print("\tchildren : \n", tabs, nil)
            pprint_children(c.children, tabs + 1)
        }
        tab_print("\t},\n", tabs, nil)
    }
    tab_print("]\n", tabs, nil)
}

func report_errors(bds []string) {
    if len(bds) > 0 {
        fmt.Printf("=== WARNING ===\nCould not access following directories:\n")
        fmt.Printf("[\n")
        for _, dir := range bds {
            fmt.Printf("\t%s,", dir)
        }
        fmt.Printf("\n]\n")
    }
}

func main() {
    // parse flags
    h := flag.Bool("hidden", false, "scan for hidden files")
    flag.Parse()

    // scans the directory the user is currently in
    // NOT the directory the executable is in!
    dir, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }

    windows := runtime.GOOS == "windows"

    if windows {
        if *h {
            flags = "/a"
        } else {
            flags = ""
        }

        info := windows_scan_dir(dir)
        report_errors(bad_dirs)

        // pprint_children([]Info{info}, 0)

        start_prompt(&info)

    } else {
        if *h {
            flags = "-la"
        } else {
            flags = "-l"
        }

        // pretty print list of size 1 consisting of returned Info struct
        //pprint_children([]Info{scan_dir(dir)}, 0)

        info := linux_scan_dir(dir)
        report_errors(bad_dirs)

        // pprint_children([]Info{info}, 0)

        start_prompt(&info)
    }


    // support for windows - ignore owner and group for now
    // refactor - some stuff (var names) are bad
    // later - generate html?, concurrency?
}


