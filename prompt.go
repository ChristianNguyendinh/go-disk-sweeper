package main

import (
    "os"
    "bufio"
    "log"
    "strings"
    "fmt"
    "strconv"
)

func show_current_info(info Info) {
    fmt.Printf("Currently Viewing: %s\n", info.name)
    for i, v := range info.children {
        fmt.Printf("\n%d.\t%s\n\tSize: %d\n", i, v.name, v.size)
    }
    fmt.Printf("\n\nPress number to go into corresponding directory\nOr back to go backwards:\n")
}

func start_prompt(info Info) {
    var curr Info = info

    fmt.Printf("Command Prompt:\n")

    show_current_info(info)

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
                if num >= int64(len(curr.children)) {
                    fmt.Printf("Invalid Number.\n")
                } else {
                    fmt.Printf("Want to go into: %s\n", curr.children[num].name)
                    curr = curr.children[num]
                    show_current_info(curr)
                }
            } else if fds[0] == "b" || fds[0] == "back" {
                fmt.Printf("Going back...\n")
                curr = *curr.parent
                show_current_info(curr)
            } else {
                fmt.Printf("Invalid: %s\n", line)
            }   
        }
    }

    fmt.Printf("DONE!\n")
}



