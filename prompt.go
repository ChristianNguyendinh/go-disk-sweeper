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

// multiple of byte sizes in binary (ex. 1 KB = 1024 B)
const (
    show_BB = 1 << (10 * iota)
    show_KB 
    show_MB
    show_GB
)

type Options struct {
    size_display int
}

var options Options
// Initialize default options
func init() {
    options.size_display = show_MB
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

func valid_format(format string) (int, bool) {
    if format == "b" || format == "B" {
        return show_BB, false
    } else if format == "kb" || format == "KB" {
        return show_KB, false
    } else if format == "mb" || format == "MB" {
        return show_MB, false
    } else if format == "gb" || format == "GB" {
        return show_GB, false
    } else {
        return 0, true
    }
}

// returns sub directories currently being shown
func show_current_info(info Info) []Info {
    fmt.Printf("\n\n______________________________________________________\n")
    fmt.Printf("Currently Viewing: %s\n", info.name)

    // sort by size, largest first
    sort.Sort(bySize(info.children))

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
        fmt.Printf("\n\t%s\n\tSize: %s\n", v.name, format_size(v.size))
    }

    fmt.Printf("\n======================\nDirectories: \n")
    for i, v := range dirs {
        fmt.Printf("\n%d.\t%s\n\tSize: %s\n", i, v.name, format_size(v.size))
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

        line := strings.TrimSpace(string(in))

        if line == "exit" {
            break
        }

        fds := strings.Fields(line)

        if len(fds) > 0 {
            if num, err := strconv.ParseInt(fds[0], 10, 64); err == nil {
                if num >= int64(len(curr_sub_dirs)) {
                    fmt.Printf("--- Invalid Number.\n")
                } else {
                    curr = curr_sub_dirs[num]
                    curr_sub_dirs = show_current_info(curr)
                }

            } else if fds[0] == "b" || fds[0] == "back" {
                curr = *curr.parent
                curr_sub_dirs = show_current_info(curr)

            } else if fds[0] == "size" && len(fds) == 2 {
                if nf, err := valid_format(fds[1]); err {
                    fmt.Printf("--- Invalid size format: %s\n", fds[1])
                } else {
                    options.size_display = nf
                    show_current_info(curr)
                    fmt.Printf("--- Now showing file sizes as %s\n", fds[1])
                }

            } else {
                fmt.Printf("--- Invalid: %s\n", line)
            }
            // help???
        }
    }

    fmt.Printf("DONE!\n")
}



