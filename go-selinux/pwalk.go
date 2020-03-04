// +build selinux,linux

package selinux

/* The code below was taken from
 * https://github.com/stretchr/powerwalk/blob/master/walker.go
 * and modified to use github.com/karrick/godirwalk
 */

import (
	"runtime"
	"sync"

	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
)

// Walk walks the file tree rooted at root, calling options.Callback
// for each file or directory in the tree, including root. All errors
// that arise visiting files and directories are filtered by Callback.
// The output is non-deterministic.
//
// For each file and directory encountered, Walk will trigger a new Go routine
// allowing you to handle each item concurrently. A maximum of twice the
// runtime.NumCPU() Callback goroutines will be called at any one time.
func Walk(root string, options *godirwalk.Options) error {
	return WalkLimit(root, options, runtime.NumCPU()*2)
}

// WalkLimit walks the file tree rooted at root, calling options.Callback
// for each file or directory in the tree, including root. All errors
// that arise visiting files and directories are filtered by Callback.
// The output is non-deterministic.
//
// For each file and directory encountered, Walk will trigger a new Go routine
// allowing you to handle each item concurrently. A maximum of limit
// Callback goroutines will be called at any one time.
func WalkLimit(root string, options *godirwalk.Options, limit int) error {
	// make sure limit is sensible
	if limit < 1 {
		return errors.Errorf("walk(%q): limit must be greater than zero", root)
	}

	files := make(chan *walkArgs, 128)
	errCh := make(chan error, 1) // get the first error, ignore others
	// save the original callback func
	callback := options.Callback

	// Start walking a tree asap
	var err error
	go func() {
		options.Callback = func(p string, info *godirwalk.Dirent) error {
			// add more files to the queue unless there's an error
			select {
			case e := <-errCh:
				close(files)
				return e
			default:
				files <- &walkArgs{path: p, info: info}
				return nil
			}
		}
		err = godirwalk.Walk(root, options)
		if err == nil {
			close(files)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(limit)
	for i := 0; i < limit; i++ {
		go func() {
			for file := range files {
				if e := callback(file.path, file.info); e != nil {
					select {
					case errCh <- e: // sent ok
					default: // buffer full
					}
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()

	return err
}

// walkArgs holds the arguments that were passed to the Walk or WalkLimit
// functions.
type walkArgs struct {
	path string
	info *godirwalk.Dirent
}
