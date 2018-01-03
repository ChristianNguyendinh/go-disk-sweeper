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

// returns sub directories currently being shown
func show_current_info(info Info) []Info {
    fmt.Printf("\n\nCurrently Viewing: %s\n", info.name)

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
        fmt.Printf("\n\t%s\n\tSize: %d\n", v.name, v.size)
    }

    fmt.Printf("\n======================\nDirectories: \n")
    for i, v := range dirs {
        fmt.Printf("\n%d.\t%s\n\tSize: %d\n", i, v.name, v.size)
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
                    fmt.Printf("Invalid Number.\n")
                } else {
                    fmt.Printf("Want to go into: %s\n", curr_sub_dirs[num].name)
                    curr = curr_sub_dirs[num]
                    curr_sub_dirs = show_current_info(curr)
                }
            } else if fds[0] == "b" || fds[0] == "back" {
                fmt.Printf("Going back...\n")
                curr = *curr.parent
                curr_sub_dirs = show_current_info(curr)
            } else {
                fmt.Printf("Invalid: %s\n", line)
            }
            // help???
        }
    }

    fmt.Printf("DONE!\n")
}



