package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/github"
	mogo "github.com/grokify/mogo/time/timeutil"
)

const (
	owner      = "ministryofjustice"
	repo       = "cloud-platform"
	escapePath = "runbooks/source/incident-log.html.md.erb"
)

func main() {
	// Parse incident-log file from github repo
	client := github.NewClient(nil)

	fileContent, _, _, err := client.Repositories.GetContents(owner, repo, escapePath, &github.RepositoryContentGetOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	if fileContent != nil {
		// Read the file line by line
		scanner := bufio.NewScanner(fileContent)
		scanner.Split(bufio.ScanLines)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		// Loop through the lines
		for i, line := range lines {
			match, err := regexp.MatchString("^## Q[1-4]", line)
			if err != nil {
				fmt.Println(err)
				return
			}
			if match {
				// parse the quarter value from the line
				quarter := line[3:10]
				// get the current quarter
				current_quarter := mogo.FormatQuarter(time.Now())
				if quarter == current_quarter {
					// goto the next line
					mttr := string(lines[i+2])
					mttr = mttr[27:]
					// take the single space between duration
					mttr = strings.ReplaceAll(mttr, " ", "")
					mttr_time, err := time.ParseDuration(mttr)
					if err != nil {
						fmt.Println(err)
						return
					}
					mttr_hours := mttr_time.Hours()
					mttr_minutes := mttr_time.Minutes()
					total_minutes := mttr_hours*60 + mttr_minutes
					fmt.Println(total_minutes)

				}
			}
		}
	}
}
