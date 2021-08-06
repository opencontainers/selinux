// +build go1.16

package pwalkdir

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func TestWalkDir(t *testing.T) {
	var count uint32
	concurrency := runtime.NumCPU() * 2

	dir, total, err := prepareTestSet(3, 2, 1)
	if err != nil {
		t.Fatalf("dataset creation failed: %v", err)
	}
	defer os.RemoveAll(dir)

	err = WalkN(dir,
		func(_ string, _ fs.DirEntry, _ error) error {
			atomic.AddUint32(&count, 1)
			return nil
		},
		concurrency)

	if err != nil {
		t.Errorf("Walk failed: %v", err)
	}
	if count != uint32(total) {
		t.Errorf("File count mismatch: found %d, expected %d", count, total)
	}

	t.Logf("concurrency: %d, files found: %d\n", concurrency, count)
}

func TestWalkDirManyErrors(t *testing.T) {
	var count uint32

	dir, total, err := prepareTestSet(3, 3, 2)
	if err != nil {
		t.Fatalf("dataset creation failed: %v", err)
	}
	defer os.RemoveAll(dir)

	max := uint32(total / 2)
	e42 := errors.New("42")
	err = Walk(dir,
		func(p string, e fs.DirEntry, _ error) error {
			if atomic.AddUint32(&count, 1) > max {
				return e42
			}
			return nil
		})
	t.Logf("found %d of %d files", count, total)

	if err == nil {
		t.Error("Walk succeeded, but error is expected")
		if count != uint32(total) {
			t.Errorf("File count mismatch: found %d, expected %d", count, total)
		}
	}
}

func makeManyDirs(prefix string, levels, dirs, files int) (count int, err error) {
	for d := 0; d < dirs; d++ {
		var dir string
		dir, err = ioutil.TempDir(prefix, "d-")
		if err != nil {
			return
		}
		count++
		for f := 0; f < files; f++ {
			fi, err := ioutil.TempFile(dir, "f-")
			if err != nil {
				return count, err
			}
			fi.Close()
			count++
		}
		if levels == 0 {
			continue
		}
		var c int
		if c, err = makeManyDirs(dir, levels-1, dirs, files); err != nil {
			return
		}
		count += c
	}

	return
}

// prepareTestSet() creates a directory tree of shallow files,
// to be used for testing or benchmarking.
//
// Total dirs: dirs^levels + dirs^(levels-1) + ... + dirs^1
// Total files: total_dirs * files
func prepareTestSet(levels, dirs, files int) (dir string, total int, err error) {
	dir, err = ioutil.TempDir(".", "pwalk-test-")
	if err != nil {
		return
	}
	total, err = makeManyDirs(dir, levels, dirs, files)
	if err != nil && total > 0 {
		_ = os.RemoveAll(dir)
		dir = ""
		total = 0
		return
	}
	total++ // this dir

	return
}

type walkerFunc func(root string, walkFn fs.WalkDirFunc) error

func genWalkN(n int) walkerFunc {
	return func(root string, walkFn fs.WalkDirFunc) error {
		return WalkN(root, walkFn, n)
	}
}

func BenchmarkWalk(b *testing.B) {
	const (
		levels = 5 // how deep
		dirs   = 3 // dirs on each levels
		files  = 8 // files on each levels
	)

	benchmarks := []struct {
		name string
		walk fs.WalkDirFunc
	}{
		{"Empty", cbEmpty},
		{"ReadFile", cbReadFile},
		{"ChownChmod", cbChownChmod},
		{"RandomSleep", cbRandomSleep},
	}

	walkers := []struct {
		name   string
		walker walkerFunc
	}{
		{"filepath.WalkDir", filepath.WalkDir},
		{"pwalkdir.Walk", Walk},
		// test WalkN with various values of N
		{"pwalkdir.Walk1", genWalkN(1)},
		{"pwalkdir.Walk2", genWalkN(2)},
		{"pwalkdir.Walk4", genWalkN(4)},
		{"pwalkdir.Walk8", genWalkN(8)},
		{"pwalkdir.Walk16", genWalkN(16)},
		{"pwalkdir.Walk32", genWalkN(32)},
		{"pwalkdir.Walk64", genWalkN(64)},
		{"pwalkdir.Walk128", genWalkN(128)},
		{"pwalkdir.Walk256", genWalkN(256)},
	}

	dir, total, err := prepareTestSet(levels, dirs, files)
	if err != nil {
		b.Fatalf("dataset creation failed: %v", err)
	}
	defer os.RemoveAll(dir)
	b.Logf("dataset: %d levels x %d dirs x %d files, total entries: %d", levels, dirs, files, total)

	for _, bm := range benchmarks {
		for _, w := range walkers {
			walker := w.walker
			walkFn := bm.walk
			// preheat
			err := w.walker(dir, bm.walk)
			if err != nil {
				b.Errorf("walk failed: %v", err)
			}
			// benchmark
			b.Run(bm.name+"/"+w.name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					err := walker(dir, walkFn)
					if err != nil {
						b.Errorf("walk failed: %v", err)
					}
				}
			})
		}
	}
}

func cbEmpty(_ string, _ fs.DirEntry, _ error) error {
	return nil
}

func cbChownChmod(path string, e fs.DirEntry, _ error) error {
	_ = os.Chown(path, 0, 0)
	mode := os.FileMode(0o644)
	if e.IsDir() {
		mode = os.FileMode(0o755)
	}
	_ = os.Chmod(path, mode)

	return nil
}

func cbReadFile(path string, e fs.DirEntry, _ error) error {
	var err error
	if e.Type().IsRegular() {
		_, err = ioutil.ReadFile(path)
	}
	return err
}

func cbRandomSleep(_ string, _ fs.DirEntry, _ error) error {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Microsecond)
	return nil
}
