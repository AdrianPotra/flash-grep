/*
Author: Adrian Potra
Version: 1.0
*/

package worker

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// our worker will be handling results, so we create the results structure - this is for when a match is found
type Result struct {
	Line       string // output of full line of text when a match is found
	LineNumber int    // which line in the file
	Path       string // so we know with which file we're working with
}

// this one contains all of the results
type Results struct {
	Inner []Result
}

// function to generate new result
func NewResult(line string, lineNum int, path string) Result {
	return Result{line, lineNum, path}
}

// function that finds data in a file - it's going to return a pointer to our Results
func FindInFile(path string, find string) *Results {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error:", err) // we report the error and then return nothing
		return nil
	}
	// we make a new results structure
	results := Results{make([]Result, 0)}

	// we will search through the file
	scanner := bufio.NewScanner(file)
	lineNum := 1

	//as long as there is data to read, we're going to search for the string and we're going to add it to our results
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), find) {
			r := NewResult(scanner.Text(), lineNum, path)
			// updating the results structure
			results.Inner = append(results.Inner, r)
		}
		lineNum += 1
	}
	//we will check if we have any results after scanner.Scan loop - if no results, we return nothing
	if len(results.Inner) == 0 {
		return nil
	} else {
		return &results
	}

}
