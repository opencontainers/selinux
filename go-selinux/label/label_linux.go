package label

import (
	"errors"
	"fmt"
	"strings"

	"github.com/opencontainers/selinux/internal/impl"
)

// Valid Label Options
var validOptions = map[string]bool{
	"disable":  true,
	"type":     true,
	"filetype": true,
	"user":     true,
	"role":     true,
	"level":    true,
}

var ErrIncompatibleLabel = errors.New("bad SELinux option: z and Z can not be used together")

// InitLabels returns the process label and file labels to be used within
// the container.  A list of options can be passed into this function to alter
// the labels.  The labels returned will include a random MCS String, that is
// guaranteed to be unique.
// If the disabled flag is passed in, the process label will not be set, but the mount label will be set
// to the container_file label with the maximum category. This label is not usable by any confined label.
func InitLabels(options []string) (plabel string, mlabel string, retErr error) {
	if !impl.GetEnabled() {
		return "", "", nil
	}
	processLabel, mountLabel := impl.ContainerLabels()
	if processLabel == "" {
		// processLabel is required; if empty, do nothing.
		return processLabel, mountLabel, nil
	}
	defer func() {
		if retErr != nil {
			impl.ReleaseLabel(mountLabel)
		}
	}()
	pcon, err := impl.NewContext(processLabel)
	if err != nil {
		return "", "", err
	}
	mcsLevel := pcon["level"]
	mcon, err := impl.NewContext(mountLabel)
	if err != nil {
		return "", "", err
	}
	for _, opt := range options {
		if opt == "disable" {
			impl.ReleaseLabel(mountLabel)
			return "", impl.PrivContainerMountLabel(), nil
		}
		k, v, ok := strings.Cut(opt, ":")
		if !ok || !validOptions[k] {
			return "", "", fmt.Errorf("bad label option %q, valid options 'disable' or \n'user, role, level, type, filetype' followed by ':' and a value", opt)
		}
		if k == "filetype" {
			mcon["type"] = v
			continue
		}
		pcon[k] = v
		if k == "level" || k == "user" {
			mcon[k] = v
		}
	}
	if p := pcon.Get(); p != processLabel {
		if pcon["level"] != mcsLevel {
			impl.ReleaseLabel(processLabel)
		}
		if err := impl.ReserveLabel(p); err != nil {
			return "", "", err
		}
		processLabel = p
	}
	mountLabel = mcon.Get()
	return processLabel, mountLabel, nil
}

// SetFileLabel modifies the "path" label to the specified file label
func SetFileLabel(path string, fileLabel string) error {
	if !impl.GetEnabled() || fileLabel == "" {
		return nil
	}
	return impl.SetFileLabel(path, fileLabel)
}

// SetFileCreateLabel tells the kernel the label for all files to be created
func SetFileCreateLabel(fileLabel string) error {
	if !impl.GetEnabled() {
		return nil
	}
	return impl.SetFSCreateLabel(fileLabel)
}

// Relabel changes the label of path and all the entries beneath the path.
// It changes the MCS label to s0 if shared is true.
// This will allow all containers to share the content.
//
// The path itself is guaranteed to be relabeled last.
func Relabel(path string, fileLabel string, shared bool) error {
	if !impl.GetEnabled() || fileLabel == "" {
		return nil
	}

	if shared {
		c, err := impl.NewContext(fileLabel)
		if err != nil {
			return err
		}

		c["level"] = "s0"
		fileLabel = c.Get()
	}
	return impl.Chcon(path, fileLabel, true)
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
