// +build selinux,linux

package selinux

/* The code below was taken from
 * https://github.com/stretchr/powerwalk/blob/master/walker.go
 * and modified to use github.com/karrick/godirwalk
 */

import (
	"sync"

	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
)

// DefaultConcurrentWalks is the default number of files that will be walked at the
// same time when the Walk function is called.
// To use a value other than this one, use the WalkLimit function.
const DefaultConcurrentWalks int = 100

// Walk walks the file tree rooted at root, calling walkFn for each file or
// directory in the tree, including root. All errors that arise visiting files
// and directories are filtered by walkFn. The output is non-deterministic.
// WalkLimit does not follow symbolic links.
//
// For each file and directory encountered, Walk will trigger a new Go routine
// allowing you to handle each item concurrently.  A maximum of DefaultConcurrentWalks
// walkFns will be called at any one time.
func Walk(root string, options *godirwalk.Options) error {
	return WalkLimit(root, options, DefaultConcurrentWalks)
}

// WalkLimit walks the file tree rooted at root, calling walkFn for each file or
// directory in the tree, including root. All errors that arise visiting files
// and directories are filtered by walkFn. The output is non-deterministic.
// WalkLimit does not follow symbolic links.
//
// For each file and directory encountered, Walk will trigger a new Go routine
// allowing you to handle each item concurrently.  A maximum of limit walkFns will
// be called at any one time.
func WalkLimit(root string, options *godirwalk.Options, limit int) error {

	// make sure limit is sensible
	if limit < 1 {
		panic("powerwalk: limit must be greater than zero.")
	}

	// filesMg is a wait group that waits for all files to
	// be processed before finishing.
	var filesWg sync.WaitGroup

	// files is a channel that receives lists of channels
	files := make(chan *walkArgs)
	kill := make(chan struct{})
	errs := make(chan error)

	for i := 0; i < limit; i++ {
		go func(i int) {
			for {
				select {
				case file, ok := <-files:
					if !ok {
						continue
					}
					if err := options.Callback(file.path, file.info); err != nil {
						errs <- err
					}
					filesWg.Done()
				case <-kill:
					return
				}
			}
		}(i)
	}

	var walkErr error

	// check for errors
	go func() {
		select {
		case walkErr = <-errs:
			close(kill)
		case <-kill:
			return
		}
	}()

	// setup a waitgroup and wait for everything to
	// be done
	var walkerWg sync.WaitGroup
	walkerWg.Add(1)

	go func() {
		opts := *options
		opts.Callback = func(p string, info *godirwalk.Dirent) error {
			select {
			case <-kill:
				close(files)
				return errors.New("kill received while walking")
			default:
				filesWg.Add(1)
				select {
				case files <- &walkArgs{path: p, info: info}:
				}
				return nil
			}
		}
		godirwalk.Walk(root, &opts)

		// everything is done
		walkerWg.Done()

	}()

	// wait for all walker calls
	walkerWg.Wait()

	if walkErr == nil {
		filesWg.Wait()
		close(kill)
	}

	return walkErr
}

// walkArgs holds the arguments that were passed to the Walk or WalkLimit
// functions.
type walkArgs struct {
	path string
	info *godirwalk.Dirent
}
