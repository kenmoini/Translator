package main

import (
	"fmt"
	"os"
	"strings"

	readingtime "github.com/begmaroman/reading-time"
)

func addReadingTime(file string) {
	// fmt.Println("Reading: ", file)
	f, err := os.ReadFile(file)
	if strings.Index(string(f), "reading_time:") > 0 {
		return
	}
	checkError(err)
	estimation := readingtime.Estimate(string(f))
	fm := strings.LastIndex(string(f), "---")
	newArt := f[:fm]
	fw, err := os.Create(file)
	checkError(err)
	defer fw.Close()
	fw.WriteString(string(newArt))
	mins := int(estimation.Duration.Minutes())
	dur := ""
	if mins > 1 {
		dur = fmt.Sprintf("reading_time: %d minutes\n", mins)
	} else if mins == 1 {
		dur = fmt.Sprintf("reading_time: %d minute\n", mins)
	} else {
		dur = fmt.Sprintf("reading_time: %d minute\n", mins)
	}
	fw.WriteString(dur)
	fw.WriteString(string(f[fm:]))
	fw.Close()
}
