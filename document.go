// A quick and easy way to review the documentation in the cloud-platform repository.
// You can execute this using the Go binary, like so:
// go run document.go --dir sources/docs

// To correct the dates simply add a --fix to the end, like so:
// go run document.go --dir sources/docs --fix
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func main() {
	var (
		dir  string
		repo string
		fix  bool
	)
	flag.StringVar(&dir, "dir", ".", "The directory to search for documents")
	flag.StringVar(&repo, "repo", "cloud-platform", "The root repository to search for documents")
	flag.BoolVar(&fix, "fix", false, "Set the out of date documents to todays date")
	flag.Parse()

	threeMonthsAgo := time.Now().AddDate(0, -3, 0).Format("2006-01-02")

	// Grab all documents needed to review.
	docs, err := allDocuments(dir)
	if err != nil {
		log.Fatalln(err)
	}

	for _, doc := range docs {
		needsReview, err := needsReview(doc, threeMonthsAgo)
		if err != nil {
			log.Printf("Unable to check if document: %s needs review %s\n", doc, err)
		}
		if needsReview {
			if fix {
				err := setLastReviewedOn(doc)
				if err != nil {
					log.Println(err)
				}
				continue
			}
			fmt.Println(doc)
		}
	}
}

func setLastReviewedOn(filePath string) error {
	// Set the last_reviewed_on date to today.
	// read the whole file at once
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	lastReviewedOn, err := lastReviewedOn(filePath)
	if err != nil {
		return errors.Wrap(err, "unable to get last_reviewed_on")
	}

	// We have to crowbar a space to make the date valid.
	todaysDate := " " + time.Now().Format("2006-01-02")

	o := bytes.Replace(b, []byte(lastReviewedOn), []byte(todaysDate), 1)
	err = ioutil.WriteFile(filePath, o, 0o644)
	if err != nil {
		return errors.Wrap(err, "unable to write file")
	}

	return nil
}

func parseTime(timeString string) (time.Time, error) {
	format := "2006-01-02"
	return time.Parse(format, strings.TrimSpace(timeString))
}

func needsReview(filePath, threeMonthsAgo string) (bool, error) {
	// Date three months ago.
	// Grab the last_reviewed_on date from the document.
	lastReviewedOn, err := lastReviewedOn(filePath)
	if err != nil {
		return false, errors.Wrap(err, "unable to get last_reviewed_on")
	}
	// Compare the last_reviewed_on date to the threeMonthsAgo date.
	lastReviewedOnTime, err := parseTime(lastReviewedOn)
	if err != nil {
		return false, fmt.Errorf("unable to parse last_reviewed_on: %s, %w", lastReviewedOn, err)
	}
	threeMonthsAgoTime, err := parseTime(threeMonthsAgo)
	if err != nil {
		return false, fmt.Errorf("unable to parse time: %s", threeMonthsAgo)
	}

	if lastReviewedOnTime.Before(threeMonthsAgoTime) {
		return true, nil
	}

	return false, nil
}

func lastReviewedOn(filePath string) (string, error) {
	// Grab the last_reviewed_on date from the document.
	// read the whole file at once
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", errors.Wrap(err, "unable to read file")
	}

	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if strings.Contains(line, "last_reviewed_on") {
			return strings.Split(line, ":")[1], nil
		}
	}
	// find the last_reviewed_on date
	return "", errors.New("unable to find last_reviewed_on")
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
		if _, ok := set[item]; ok {
			return true
		}
	}
	return false
}

func allDocuments(dir string) ([]string, error) {
	// documents is a map of document names and dates to review them by.
	ignored := []string{
		"incident-log.html.md.erb",
		"index.html.md.erb",
		"custom.erb",
	}

	var documents []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".erb" && !contains(ignored, info.Name()) {
			documents = append(documents, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return documents, nil
}
