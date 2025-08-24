package main

import (
	"fmt"
	"os"
	"time"

	"github.com/JuanMartinCoder/portanalyzer/internal/portscanner"
	"github.com/JuanMartinCoder/portanalyzer/internal/utils"
	"golang.org/x/sync/semaphore"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Correct Formart: go run ./cmd [ip-addr]")
		os.Exit(1)
	}
	targetedIp := os.Args[1]
	timeNow := time.Now()
	fmt.Printf("Starting port scanner at %s \n\n", timeNow.Format(time.RFC3339))

	portscanner := portscanner.NewPortScanner(targetedIp, semaphore.NewWeighted(utils.Ulimit()))
	portscanner.Scan(65535, 500*time.Millisecond)
	portscanner.ShowOpenPorts()
	elapsed := time.Since(timeNow)

	fmt.Printf("\ndone: 1 IP address (%s) scanned in %s \n", targetedIp, elapsed.String())
}
