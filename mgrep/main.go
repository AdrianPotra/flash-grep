/*
Author: Adrian Potra
Version 1.0
*/

package main

import (
	"fmt"
	"mgrep/worker"
	"mgrep/worklist"
	"os"
	"path/filepath"
	"sync"

	"github.com/alexflint/go-arg"
)

// function to locate all of the directories and files - it's going to take a pointer to the worklist along with an initial search path
func discoverDirs(wl *worklist.Worklist, path string) {
	entries, err := os.ReadDir(path) // gives us a listing with everything in the directory
	if err != nil {
		fmt.Println("ReadDir Error:", err)
		return
	}

	for _, entry := range entries {
		//fmt.Println("entry value in range is", entry)
		if entry.Type().IsDir() { // if we run into a directory, we need to recurse into it
			nextPath := filepath.Join(path, entry.Name()) //this is constructing the path with the new directory - Example: if path given was "a" and it found also "b", the file path should be "a/b"
			//fmt.Println("FilePath constructed when doing filepath join and taking just the name", nextPath)
			// now we are going to recurse into the directory
			discoverDirs(wl, nextPath)
		} else { // if it's not, then add it to the worklist
			//file := filepath.Join(path, entry.Name())
			//fmt.Println("file in directory recursion is: ", file)
			wl.Add(worklist.NewJob(filepath.Join(path, entry.Name())))

		}
	}
}

// using the github go-arg parser - will attach some annotations - the package reads the annotations - since we will have one of them required, if we don't provide it, the program should abort
var args struct {
	SearchTerm string `arg:"positional,required"`
	SearchDir  string `arg:"positional"`
}

func main() {

	// we need to parse the command line args
	arg.MustParse(&args)

	// create a work group for workers
	var workersWg sync.WaitGroup

	// creating new worklist with size 100
	wl := worklist.New(100)

	//results channel where the workers will put their work in
	results := make(chan worker.Result, 100)

	//setting number of workers
	numWorkers := 10
	// HERE TO MAKE ADDITION TO INPUT IN OWN TERMINAL - take into consideration line 71 and 84 - the args. - should input  plus
	// adding instructions

	// launching go routine for discovering directories
	workersWg.Add(1)
	go func() {
		defer workersWg.Done()
		discoverDirs(&wl, args.SearchDir) // we pull our search directory
		wl.Finalize(numWorkers)           // empty pads in the end to indicate that workers should end

	}()

	// spawning workers loop
	for i := 0; i < numWorkers; i++ {
		workersWg.Add(1)
		go func() {
			defer workersWg.Done()
			for { //in this infinite loop we're going to use the find in file func that we created for the worker
				workEntry := wl.Next() // pulling out a work entry by accessing the next func, which was populated by the worklist go routine
				if workEntry.Path != "" {
					workerResult := worker.FindInFile(workEntry.Path, args.SearchTerm) // obtaining results
					if workerResult != nil {                                           // if we do get a result, we loop through it/them and for each result we're going to put it into the results channel
						for _, r := range workerResult.Inner {
							results <- r
						}

					}
				} else {
					return // empty means no work to be done
				}

			}

		}()
	}

	// wait group channel for blocking
	blockWorkersWg := make(chan struct{})
	go func() { // this go routine will sit down blocking while the workes finish
		workersWg.Wait()
		close(blockWorkersWg)
	}()

	// creating a channel for displaying results while the workers are working
	var displayWg sync.WaitGroup
	displayWg.Add(1)

	go func() {
		for {
			select {
			case r := <-results:
				fmt.Printf("%v[%v]:%v\n", r.Path, r.LineNumber, r.Line) // %v to get the default values from the path line
			case <-blockWorkersWg:
				if len(results) == 0 { // if there are no results, then this will terminate - if there are results, then it goes back to results case to print out
					displayWg.Done()
					return
				}
			}
		}
	}()

	displayWg.Wait()
}
