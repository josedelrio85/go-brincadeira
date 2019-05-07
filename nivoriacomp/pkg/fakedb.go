package nivoriacomp

import "sync"

// FakeDb is a struct used to test Db functionality with fake methods.
type FakeDb struct {
	OpenFunc              func() error
	OpenCalls             int
	BatchInsertFunc       func([]Xmlstruct) error
	BatchInsertCalls      int
	SelectForRequestFunc  func([]Inputdata) error
	SelectForRequestCalls int

	sync.Mutex
}

// Open is a method to test Open function
func (f *FakeDb) Open() error {
	f.Lock()
	defer f.Unlock()
	f.OpenCalls++
	return f.OpenFunc()
}

// BatchInsert is a method to test BatchInsert function
func (f *FakeDb) BatchInsert(rows []Xmlstruct) error {
	f.Lock()
	defer f.Unlock()
	f.BatchInsertCalls++
	return f.BatchInsertFunc(rows)
}

// SelectForRequest is a method to test BatchInsert function
func (f *FakeDb) SelectForRequest(inputdata []Inputdata) error {
	f.Lock()
	defer f.Unlock()
	f.SelectForRequestCalls++
	return f.SelectForRequestFunc(inputdata)
}
