package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	filePath     = "displaynumberofevents/events.log"
	suffix       = "NOK"
	parseLayout  = "[2006-01-02 15:04:05]"
	outputLayout = "[2006-01-02 15:04]"
)

type output struct {
	minuteTime time.Time
	count      int
}

func openFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func main() {
	logFile := openFile(filePath)
	defer logFile.Close()

	scanner := bufio.NewScanner(logFile)

	result := make([]output, 1)

	for scanner.Scan() {

		// skip empty
		if len(scanner.Text()) < len(parseLayout) {
			continue
		}

		//	parse string date to golang time,
		//  also stop runtime and print error if line did not validate format
		t, err := time.Parse(parseLayout, scanner.Text()[:len(parseLayout)])
		if err != nil {
			log.Fatal(err)
		}

		tMin := t.Truncate(time.Minute)
		if result[len(result)-1].minuteTime != tMin {
			result = append(result, output{
				minuteTime: tMin,
				count:      0,
			})
		}

		if strings.HasSuffix(scanner.Text(), suffix) {
			result[len(result)-1].count++
		}

	}

	for _, v := range result[1:] {
		fmt.Println(v.minuteTime.Format(outputLayout), v.count)
	}
}
