package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func readLines(path string) ([]string, error) {
	var lines []string

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func convertToRaw(data string) int {
	s := strings.Split(data, " ")
	hours := 0
	minutes := 0

	for _, time := range s {
		if strings.Contains(time, "h") {
			time = strings.Replace(time, "h", "", -1)
			i, err := strconv.Atoi(time)
			if err != nil {
				log.Fatalf("contains h: %s", err)
			}
			hours = convertHoursToMinutes(i)
		}
		if strings.Contains(time, "m") {
			time = strings.Replace(time, "m", "", -1)
			i, err := strconv.Atoi(time)
			if err != nil {
				log.Fatalf("contains m: %s", err)
			}
			minutes = i
		}
	}

	return hours + minutes
}

func convertHoursToMinutes(i int) int {
	return i * 60
}

func main() {
	var data []string
	var str strings.Builder

	lines, err := readLines("../../runbooks/source/incident-log.html.md.erb")
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	for _, line := range lines {
		str.WriteString(line)

		re := regexp.MustCompile(`---`)
		match := re.FindString(line)

		if match != "" {
			data = append(data, str.String())
			str.Reset()
		}
	}

	for _, newline := range data {
		reTitle := regexp.MustCompile(`(?U)## .\d \d* \(.*\)`)

		title := reTitle.FindString((newline))

		if title != "" {
			fmt.Printf("%s\n", title)
		}

		re := regexp.MustCompile(`\*\*Time to repair\*\*: (\d*. \d*.|\d*.)`)
		timeToRepair := 0
		count := 0
		for _, regmatch := range re.FindAllString(newline, -1) {
			t := strings.Replace(regmatch, "**Time to repair**: ", "", -1)
			timeToRepairTemp := convertToRaw(t)
			timeToRepair = timeToRepair + timeToRepairTemp
			count += 1
		}

		re2 := regexp.MustCompile(`\*\*Time to resolve\*\*: (\d*. \d*.|\d*.)`)
		timeToResolve := 0
		resolveCount := 0
		for _, resolveMatch := range re2.FindAllString(newline, -1) {
			t := strings.Replace(resolveMatch, "**Time to resolve**: ", "", -1)
			timeToResolveTemp := convertToRaw(t)
			timeToResolve = timeToResolve + timeToResolveTemp
			resolveCount += 1
		}

		if count != 0 {
			meanTimeToRepair := timeToRepair / count
			d := time.Duration(meanTimeToRepair) * time.Minute
			hours := int(d.Hours())
			minutes := int(d.Minutes()) % 60
			fmt.Printf("Incidents:%2d\n", count)
			fmt.Printf("Mean time to repair: %2dh %02dm\n", hours, minutes)
		}

		if resolveCount != 0 {
			meantTimeToResolve := timeToResolve / resolveCount
			d := time.Duration(meantTimeToResolve) * time.Minute
			hours := int(d.Hours())
			minutes := int(d.Minutes()) % 60
			fmt.Printf("Mean time to resolve: %2dh %02dm", hours, minutes)
			fmt.Println("\n")
		}

	}
}
