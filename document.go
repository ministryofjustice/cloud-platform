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
		dir string
		fix bool
	)
	flag.StringVar(&dir, "dir", "runbooks/source", "The directory to search for documents")
	flag.BoolVar(&fix, "fix", false, "Set the out of date documents to todays date")
	flag.Parse()

	threeMonthsAgo := time.Now().AddDate(0, -3, 0).Format("2006-01-02")

	err := checkWorkingDir(dir)
	if err != nil {
		// Hard error if the user executes from the wrong directory.
		log.Fatalln(err)
	}

	// Grab all documents needed to review.
	docs, err := allDocuments(dir)
	if err != nil {
		log.Fatalln(err)
	}

	if fix {
		for _, doc := range docs {
			err := setLastReviewedOn(doc)
			if err != nil {
				log.Fatalln(err)
			}
		}
	} else {
		for _, doc := range docs {
			needsReview, err := needsReview(doc, threeMonthsAgo)
			if err != nil {
				log.Printf("Unable to check if document: %s needs review %s\n", doc, err)
			}
			if needsReview {
				fmt.Println(doc)
			}
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

	o := bytes.Replace(b, []byte(lastReviewedOn), []byte(time.Now().Format("2006-01-02")), 1)
	err = ioutil.WriteFile(filePath, o, 0644)
	if err != nil {
		return errors.Wrap(err, "unable to write file")
	}

	return nil
}

func checkWorkingDir(workingDir string) error {
	// Am I in the `cloud-platform` directory?
	const expectedDir = "cloud-platform"
	repo, err := findTopLevelGitDir("./")
	if err != nil {
		return fmt.Errorf("unable to get the repository name, %e", err)
	}

	if !strings.Contains(repo, expectedDir) {
		return fmt.Errorf("you must execute in the %s repoa, %e", expectedDir, err)
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

func findTopLevelGitDir(workingDir string) (string, error) {
	dir, err := filepath.Abs(workingDir)
	if err != nil {
		return "", errors.Wrap(err, "invalid working dir")
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("no git repository found")
		}
		dir = parent
	}
}
