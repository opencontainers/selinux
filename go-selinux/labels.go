package selinux

import (
	"fmt"
	"strings"
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

// InitLabels returns the process label and file labels to be used within
// the container.  A list of options can be passed into this function to alter
// the labels.  The labels returned will include a random MCS String, that is
// guaranteed to be unique.
// If the disabled flag is passed in, the process label will not be set, but the mount label will be set
// to the container_file label with the maximum category. This label is not usable by any confined label.
func InitLabels(options []string) (plabel string, mlabel string, retErr error) {
	if !GetEnabled() {
		return "", "", nil
	}
	processLabel, mountLabel, err := containerLabels()
	if err != nil {
		return "", "", err
	}
	if processLabel == "" {
		// processLabel is required; if empty, do nothing.
		return processLabel, mountLabel, nil
	}
	defer func() {
		if retErr != nil {
			ReleaseLabel(mountLabel)
		}
	}()
	pcon, err := NewContext(processLabel)
	if err != nil {
		return "", "", err
	}
	mcsLevel := pcon["level"]
	mcon, err := NewContext(mountLabel)
	if err != nil {
		return "", "", err
	}
	for _, opt := range options {
		if opt == "disable" {
			ReleaseLabel(mountLabel)
			return "", PrivContainerMountLabel(), nil
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
			ReleaseLabel(processLabel)
		}
		ReserveLabel(p)
		processLabel = p
	}
	mountLabel = mcon.Get()
	return processLabel, mountLabel, nil
}
