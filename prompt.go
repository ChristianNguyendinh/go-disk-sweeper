// Functions for displaying the prompt for showing and traversing gathered info
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

func valid_size(f string) (int, bool) {
    switch f {
    case "b":
        return show_BB, false
    case "kb":
        return show_KB, false
    case "mb":
        return show_MB, false
    case "gb":
        return show_GB, false
    default:
        return 0, true
    }
}

func show_size(f int) string {
    switch f {
    case show_BB:
        return "Bytes (b)"
    case show_KB:
        return "Kilobytes (kb)"
    case show_MB:
        return "Megabytes (mb)"
    case show_GB:
        return "Gigabytes (gb)"
    default:
        return "invalid"
    }
}

func valid_sort(s string) (int, bool) {
    switch s {
    case "size":
        return sort_size, false
    case "name":
        return sort_name, false
    case "owner":
        return sort_owner, false
    case "group":
        return sort_group, false
    default:
        return 0, true
    }
}

func show_sort(s int) string {
    switch s {
    case sort_size:
        return "size (largest on top)"
    case sort_name:
        return "name"
    case sort_owner:
        return "owner"
    case sort_group:
        return "group"
    default:
        return "invalid"
    }
}

// returns sub directories currently being shown
func show_current_info(info *Info) []*Info {
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

    dirs := []*Info{}
    files := []*Info{}
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

func start_prompt(info *Info) {
    home := info
    curr := info
    curr_sub_dirs := show_current_info(curr)

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
                    // WHY DO I NEED TO DO THIS?!?!?!?!?!?!?!?!?!?!??!?
                    // fixed by having the children be array of ptrs instead of structs
                    // i still dont know wtf why
                    // curr.parent = par
                }

            } else if fds[0] == "b" || fds[0] == "back" {
                if curr.parent != nil {
                    curr = curr.parent
                    curr_sub_dirs = show_current_info(curr)
                } else {
                    fmt.Printf("Cannot go back any more.\n")
                }

            } else if fds[0] == "size" && len(fds) <= 2 {
                if len(fds) == 1 {
                    fmt.Printf("Currently showing file sizes in: %s\n", show_size(options.size_display))
                } else {
                    if nf, err := valid_size(fds[1]); err {
                        fmt.Printf("!!! Invalid size format: %s\n", fds[1])
                    } else {
                        options.size_display = nf
                        curr_sub_dirs = show_current_info(curr)
                        fmt.Printf("--- Now showing file sizes as %s\n", fds[1])
                    }
                }

            } else if fds[0] == "sort" && len(fds) <= 2 {
                if len(fds) == 1 {
                    fmt.Printf("Currently sorting by: %s\n", show_sort(options.sorting))
                } else {
                    if sf, err := valid_sort(fds[1]); err {
                        fmt.Printf("!!! Invalid sort format: %s\n", fds[1])
                    } else {
                        options.sorting = sf
                        curr_sub_dirs = show_current_info(curr)
                        fmt.Printf("--- Now sorting by %s\n", fds[1])
                    }
                }

            } else if fds[0] == "verbose" && len(fds) <= 2 {
                if len(fds) == 1 {
                    if options.verbose {
                        fmt.Printf("verbose currently on\n")
                    } else {
                        fmt.Printf("verbose currently off\n")
                    }
                } else {
                    if fds[1] == "on" {
                        options.verbose = true
                        curr_sub_dirs = show_current_info(curr)
                        fmt.Printf("--- Showing all info\n")
                    } else if fds[1] == "off"{
                        options.verbose = false
                        curr_sub_dirs = show_current_info(curr)
                        fmt.Printf("--- Only showing name and size\n")
                    } else {
                        fmt.Printf("!!! Invalid sort format: %s\n", fds[1])
                    }
                }

            } else if len(fds) == 1 && (fds[0] == "home" || fds[0] == "~") {
                curr = home
                curr_sub_dirs = show_current_info(curr)

            } else if fds[0] == "help" {
                fmt.Printf("\n--- List of Commands:\n")
                fmt.Printf("\t<number> - Go into the corresponding directory\n\n")
                fmt.Printf("\tb\n\tback - Go into the corresponding directory\n\n")
                fmt.Printf("\tsize (b|kb|mb|gb) - change display size\n\n")
                fmt.Printf("\tsize - show current display size\n\n")
                fmt.Printf("\tsort (size|name|owner|group) - change order which files/dirs are displayed\n\n")
                fmt.Printf("\tsort - show current sorting display method\n\n")
                fmt.Printf("\tverbose (on|off) - for each file/dir, on shows owner and group, off does not\n\n")
                fmt.Printf("\tverbose - show whether verbose is currently on or not\n\n")
                fmt.Printf("\thelp - show command information\n\n")
                fmt.Printf("\texit - exit the program\n\n")


            } else {
                fmt.Printf("!!! Invalid Command: %s\n", line)
            }
        }
    }

    fmt.Printf("\nexiting...\n\n")
}



