/*
Author: Adrian Potra
Version 1.0
*/

package worklist

type Entry struct {
	Path string
}

// worklist is a channel and is going to hold the jobs of type Entry
type Worklist struct {
	jobs chan Entry
}

//receiver function to add a job
func (w *Worklist) Add(work Entry) {
	w.jobs <- work
}

//receiver function to get the job out of the channel
func (w *Worklist) Next() Entry {
	j := <-w.jobs
	return j
}

//function to generate new worklist - with creation of a buffered channel
func New(bufSize int) Worklist {
	return Worklist{make(chan Entry, bufSize)}
}

//function to generate new job
func NewJob(path string) Entry {
	return Entry{path}
}

// function to generate empty jobs - this will be signalling to the workers that it's time for them to quit (no more files available)
// the for loop is to terminate every single active worker and to signal to them all to terminate
func (w *Worklist) Finalize(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		w.Add(Entry{""})
	}
}
