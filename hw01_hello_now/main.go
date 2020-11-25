package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

const NTPServer = "ntp3.stratum2.ru"
const timeFmt = "2006-01-02 03:04:05 -0700 MST"

func main() {
	currentTime := time.Now().UTC()
	exactTimeNtp, err := ntp.Time(NTPServer)
	if err != nil {
		log.Fatal(err)
	}
	exactTime := exactTimeNtp.UTC()
	fmt.Printf("current time: %s\nexact time: %s\n", currentTime.Format(timeFmt), exactTime.Format(timeFmt))
}
