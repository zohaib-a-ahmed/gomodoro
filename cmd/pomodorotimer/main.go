package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "strings"
    "time"
)

func main() {
    // Define program arguments
    workDuration := flag.Int("work", 25, "Specify work length (minutes)")
    breakDuration := flag.Int("break", 5, "Specify break length (minutes)")

    flag.Parse() // Parse the command-line flags

    controlChan := make(chan string)
    go handleUserInput(controlChan)

    var title string
    title = `
        ____   _____    _   _   _____    ____     _____    _____    _____
       |    | |     |  | | | | |     |  |    |   |     |  |     |  |     |
       |____| |     |  | |_| | |     |  |     |  |     |  |_____|  |     |
       |      |_____|  |     | |_____|  |____|   |_____|  |   \    |_____|
    `

    fmt.Printf(title)
    fmt.Println("")
    fmt.Printf("Commands: 'start' to begin, 'pause' to pause, 'resume' to resume, 'quit' to quit\n")
    for command := range controlChan {
        switch command {
        case "start":
            fmt.Println("Starting work timer...")
            runTimer(time.Duration(*workDuration)*time.Minute, controlChan)
            fmt.Println("Work session complete! Starting break timer automatically...")
            runTimer(time.Duration(*breakDuration)*time.Minute, controlChan)
            fmt.Println("Break session complete! Timer ended.")
            return
        case "quit":
            fmt.Println("Quitting the timer.")
            return
        }
    }
}

func handleUserInput(controlChan chan<- string) {
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        controlChan <- line
    }
    if scanner.Err() != nil {
        fmt.Fprintf(os.Stderr, "Error reading from stdin: %s\n", scanner.Err())
    }
}

func runTimer(duration time.Duration, controlChan <-chan string) {
    const progressBarWidth = 50
    startTime := time.Now()
    endTime := startTime.Add(duration)

    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()

    pause := false
    for {
        select {
        case <-ticker.C:
            if !pause {
                //elapsed := time.Since(startTime)
                remaining := endTime.Sub(time.Now())

                if remaining <= 0 {
                    fmt.Println("\nTimer completed!")
                    return
                }

                percentage := 1 - remaining.Seconds()/duration.Seconds()
                fill := int(percentage * float64(progressBarWidth))

                progressBar := fmt.Sprintf("\r[%s%s] %2d%% Remaining: %v", strings.Repeat("#", fill), strings.Repeat("-", progressBarWidth-fill), int(percentage*100), remaining.Truncate(time.Second))
                fmt.Print(progressBar)
            }
        case cmd := <-controlChan:
            if cmd == "pause" {
                pause = true
                ticker.Stop()
            } else if cmd == "resume" {
                pause = false
                ticker = time.NewTicker(100 * time.Millisecond)
                startTime = time.Now().Add(endTime.Sub(time.Now()))
            } else if cmd == "quit" {
                fmt.Println("\nQuitting early...")
                return
            }
        }
    }
}
