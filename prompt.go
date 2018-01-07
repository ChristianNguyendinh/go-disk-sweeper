package main

import (
    "os"
    "bufio"
    "log"
    "strings"
    "fmt"
    "strconv"
    "sort"
)

const (
    // multiple of byte sizes in binary (ex. 1 KB = 1024 B)
    show_BB = 1 << (10 * iota)
    show_KB 
    show_MB
    show_GB

    sort_size = iota
    sort_name
    sort_owner
    sort_group
)

type Options struct {
    size_display    int
    sorting         int
    verbose         bool
}

var options Options
// Initialize default options
func init() {
    options.size_display = show_BB
    options.sorting = sort_size
    options.verbose = true
}

func format_size(size uint64) string {
    switch options.size_display {
    case show_KB:
        return fmt.Sprintf("%.3fKB", float64(size) / show_KB)
    case show_MB:
        return fmt.Sprintf("%.3fMB", float64(size) / show_MB)
    case show_GB:
        return fmt.Sprintf("%.3fGB", float64(size) / show_GB)
    default:
        return fmt.Sprintf("%dB", size)
    }
}

func valid_format(f string) (int, bool) {
    if f == "b" {
        return show_BB, false
    } else if f == "kb" {
        return show_KB, false
    } else if f == "mb" {
        return show_MB, false
    } else if f == "gb" {
        return show_GB, false
    } else {
        return 0, true
    }
}

func valid_sort(s string) (int, bool) {
    if s == "size" {
        return sort_size, false
    } else if s == "name" {
        return sort_name, false
    } else if s == "owner" {
        return sort_owner, false
    } else if s == "group" {
        return sort_group, false
    } else {
        return 0, true
    }
}

// returns sub directories currently being shown
func show_current_info(info Info) []Info {
    fmt.Printf("\n\n______________________________________________________\n")
    fmt.Printf("Currently Viewing: %s\n", info.name)

    // sort results
    switch options.sorting {
    case sort_size:
        sort.Sort(bySize(info.children))
    case sort_name:
        sort.Sort(byName(info.children))
    case sort_owner:
        sort.Sort(byOwner(info.children))
    case sort_group:
        sort.Sort(byGroup(info.children))
    }

    dirs := []Info{}
    files := []Info{}
    // separate files and dirs, keeping sorted order
    for _, v := range info.children {
        if v.directory {
            dirs = append(dirs, v)
        } else {
            files = append(files, v)
        }
    }

    fmt.Printf("\n======================\nFiles: \n")
    for _, v := range files {
        fmt.Printf("\n-\t%s\n\tSize: %s\n", v.name, format_size(v.size))
        if options.verbose {
            fmt.Printf("\tOwner: %s\n\tGroup: %s\n", v.owner, v.group)
        }
    }

    fmt.Printf("\n======================\nDirectories: \n")
    for i, v := range dirs {
        fmt.Printf("\n%d.\t%s\n\tSize: %s\n", i, v.name, format_size(v.size))
        if options.verbose {
            fmt.Printf("\tOwner: %s\n\tGroup: %s\n", v.owner, v.group)
        }
    }

    fmt.Printf("\n\nPress number to go into corresponding directory\nOr back to go backwards:\n")

    return dirs
}

func start_prompt(info Info) {
    curr := info
    curr_sub_dirs := show_current_info(info)

    for {
        bio := bufio.NewReader(os.Stdin)
        fmt.Printf(" > ")
        in, _, err := bio.ReadLine()
        if err != nil {
            log.Fatal(err)
        }

        line := strings.ToLower(strings.TrimSpace(string(in)))

        if line == "exit" {
            break
        }

        fds := strings.Fields(line)

        if len(fds) > 0 {
            if num, err := strconv.ParseInt(fds[0], 10, 64); err == nil {
                if num >= int64(len(curr_sub_dirs)) {
                    fmt.Printf("!!! Invalid Number.\n")
                } else {
                    curr = curr_sub_dirs[num]
                    curr_sub_dirs = show_current_info(curr)
                }

            } else if fds[0] == "b" || fds[0] == "back" {
                curr = *curr.parent
                curr_sub_dirs = show_current_info(curr)

            } else if fds[0] == "size" && len(fds) == 2 {
                if nf, err := valid_format(fds[1]); err {
                    fmt.Printf("!!! Invalid size format: %s\n", fds[1])
                } else {
                    options.size_display = nf
                    show_current_info(curr)
                    fmt.Printf("--- Now showing file sizes as %s\n", fds[1])
                }

            } else if fds[0] == "sort" && len(fds) == 2 {
                if sf, err := valid_sort(fds[1]); err {
                    fmt.Printf("!!! Invalid sort format: %s\n", fds[1])
                } else {
                    options.sorting = sf
                    show_current_info(curr)
                    fmt.Printf("--- Now sorting by %s\n", fds[1])
                }

            } else if fds[0] == "verbose" && len(fds) == 2 {
                if fds[1] == "on" {
                    options.verbose = true
                    show_current_info(curr)
                    fmt.Printf("--- Showing all info\n")
                } else if fds[1] == "off"{
                    options.verbose = false
                    show_current_info(curr)
                    fmt.Printf("--- Only showing name and size\n")
                } else {
                    fmt.Printf("!!! Invalid sort format: %s\n", fds[1])
                }

            } else if fds[0] == "help" {
                fmt.Printf("\n--- List of Commands:\n")
                fmt.Printf("\t<number> - Go into the corresponding directory\n\n")
                fmt.Printf("\tb\n\tback - Go into the corresponding directory\n\n")
                fmt.Printf("\tsize (b|kb|mb|gb) - change display size\n\n")
                fmt.Printf("\tsort (size|name|owner|group) - change order which files/dirs are displayed\n\n")
                fmt.Printf("\tverbose (on|off) - for each file/dir, on shows owner and group, off does not\n\n")
                fmt.Printf("\thelp - show command information\n\n")
                fmt.Printf("\texit - exit the program\n\n")


            } else {
                fmt.Printf("!!! Invalid Command: %s\n", line)
            }
        }
    }

    fmt.Printf("DONE!\n")
}



