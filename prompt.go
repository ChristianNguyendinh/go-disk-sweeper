package main

import (
    "os"
    "bufio"
    "log"
    "strings"
    "fmt"
)

func show_current_info(info Info) {
    fmt.Printf("Currently Viewing: %s\n", info.name)
    for i, v := range info.children {
        fmt.Printf("\n%d.\t%s\n\tSize: %d\n", i, v.name, v.size)
    }
    fmt.Printf("\n\nPress number to go into corresponding directory\nOr back to go backwards:\n")
}

func start_prompt(info Info) {
    // var curr Info = info
    // var prev Info = info

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

        // fds = strings.Fields(line)

        // switch fds[0] {
        // case "":
        // }

        fmt.Printf("Command Ran: %s\n", line)
    }

    fmt.Printf("DONE!\n")
}



