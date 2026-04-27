package label

import (
	"errors"
	"strings"

	"github.com/opencontainers/selinux/go-selinux"
)

var ErrIncompatibleLabel = errors.New("bad SELinux option: z and Z can not be used together")

// InitLabels returns the process label and file labels to be used within
// the container.  A list of options can be passed into this function to alter
// the labels.  The labels returned will include a random MCS String, that is
// guaranteed to be unique.
// If the disabled flag is passed in, the process label will not be set, but the mount label will be set
// to the container_file label with the maximum category. This label is not usable by any confined label.
func InitLabels(options []string) (plabel string, mlabel string, retErr error) {
	return selinux.InitLabels(options)
}

// SetFileLabel modifies the "path" label to the specified file label
func SetFileLabel(path string, fileLabel string) error {
	if !selinux.GetEnabled() || fileLabel == "" {
		return nil
	}
	return selinux.SetFileLabel(path, fileLabel)
}

// SetFileCreateLabel tells the kernel the label for all files to be created
func SetFileCreateLabel(fileLabel string) error {
	if !selinux.GetEnabled() {
		return nil
	}
	return selinux.SetFSCreateLabel(fileLabel)
}

// Relabel changes the label of path and all the entries beneath the path.
// It changes the MCS label to s0 if shared is true.
// This will allow all containers to share the content.
//
// The path itself is guaranteed to be relabeled last.
func Relabel(path string, fileLabel string, shared bool) error {
	if !selinux.GetEnabled() || fileLabel == "" {
		return nil
	}

	if shared {
		c, err := selinux.NewContext(fileLabel)
		if err != nil {
			return err
		}

		c["level"] = "s0"
		fileLabel = c.Get()
	}
	return selinux.Chcon(path, fileLabel, true)
}

// Validate checks that the label does not include unexpected options
func Validate(label string) error {
	if strings.Contains(label, "z") && strings.Contains(label, "Z") {
		return ErrIncompatibleLabel
	}
	return nil
}

// RelabelNeeded checks whether the user requested a relabel
func RelabelNeeded(label string) bool {
	return strings.Contains(label, "z") || strings.Contains(label, "Z")
}

// IsShared checks that the label includes a "shared" mark
func IsShared(label string) bool {
	return strings.Contains(label, "z")
}
