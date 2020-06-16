// +build !selinux !linux

package selinux

// setDisabled disables selinux support for the package
func setDisabled() {
}

// getEnabled returns whether selinux is currently enabled.
func getEnabled() bool {
	return false
}

// classIndex returns the int index for an object class in the loaded policy, or -1 and an error
func classIndex(class string) (int, error) {
	return -1, nil
}

// SetFileLabel sets the SELinux label for this path or returns an error.
func setFileLabel(fpath string, label string) error {
	return nil
}

// fileLabel returns the SELinux label for this path or returns an error.
func fileLabel(fpath string) (string, error) {
	return "", nil
}

/*
setFSCreateLabel tells kernel the label to create all file system objects
created by this task. Setting label="" to return to default.
*/
func setFSCreateLabel(label string) error {
	return nil
}

/*
fsCreateLabel returns the default label the kernel which the kernel is using
for file system objects created by this task. "" indicates default.
*/
func fsCreateLabel() (string, error) {
	return "", nil
}

// currentLabel returns the SELinux label of the current process thread, or an error.
func currentLabel() (string, error) {
	return "", nil
}

// PidLabel returns the SELinux label of the given pid, or an error.
func pidLabel(pid int) (string, error) {
	return "", nil
}

/*
execLabel returns the SELinux label that the kernel will use for any programs
that are executed by the current process thread, or an error.
*/
func execLabel() (string, error) {
	return "", nil
}

/*
canonicalizeContext takes a context string and writes it to the kernel
the function then returns the context that the kernel will use.  This function
can be used to see if two contexts are equivalent
*/
func canonicalizeContext(val string) (string, error) {
	return "", nil
}

/*
computeCreateContext requests the type transition from source to target for class  from the kernel.
*/
func computeCreateContext(source string, target string, class string) (string, error) {
	return "", nil
}

// calculateGlbLub computes the glb (greatest lower bound) and lub (least upper bound)
// of a source and target range.
// The glblub is calculated as the greater of the low sensitivities and
// the lower of the high sensitivities and the and of each category bitmap.
func calculateGlbLub(sourceRange, targetRange string) (string, error) {
	return "", nil
}

/*
setExecLabel sets the SELinux label that the kernel will use for any programs
that are executed by the current process thread, or an error.
*/
func setExecLabel(label string) error {
	return nil
}

/*
setTaskLabel sets the SELinux label for the current thread, or an error.
This requires the dyntransition permission.
*/
func setTaskLabel(label string) error {
	return nil
}

/*
setSocketLabel sets the SELinux label that the kernel will use for any programs
that are executed by the current process thread, or an error.
*/
func setSocketLabel(label string) error {
	return nil
}

// socketLabel retrieves the current socket label setting
func socketLabel() (string, error) {
	return "", nil
}

// peerLabel retrieves the label of the client on the other side of a socket
func peerLabel(fd uintptr) (string, error) {
	return "", nil
}

// setKeyLabel takes a process label and tells the kernel to assign the
// label to the next kernel keyring that gets created
func setKeyLabel(label string) error {
	return nil
}

// keyLabel retrieves the current kernel keyring label setting
func keyLabel() (string, error) {
	return "", nil
}

// Get returns the Context as a string
func (c Context) get() string {
	return ""
}

// newContext creates a new Context struct from the specified label
func newContext(label string) (Context, error) {
	c := make(Context)
	return c, nil
}

// clearLabels clears all reserved MLS/MCS levels
func clearLabels() {
}

// reserveLabel reserves the MLS/MCS level component of the specified label
func reserveLabel(label string) {
}

// enforceMode returns the current SELinux mode Enforcing, Permissive, Disabled
func enforceMode() int {
	return Disabled
}

/*
setEnforceMode sets the current SELinux mode Enforcing, Permissive.
Disabled is not valid, since this needs to be set at boot time.
*/
func setEnforceMode(mode int) error {
	return nil
}

/*
defaultEnforceMode returns the systems default SELinux mode Enforcing,
Permissive or Disabled. Note this is is just the default at boot time.
EnforceMode tells you the systems current mode.
*/
func defaultEnforceMode() int {
	return Disabled
}

/*
releaseLabel will unreserve the MLS/MCS Level field of the specified label.
Allowing it to be used by another process.
*/
func releaseLabel(label string) {
}

// roFileLabel returns the specified SELinux readonly file label
func roFileLabel() string {
	return ""
}

// kvmContainerLabels returns the default processLabel and mountLabel to be used
// for kvm containers by the calling process.
func kvmContainerLabels() (string, string) {
	return "", ""
}

// initContainerLabels returns the default processLabel and file labels to be
// used for containers running an init system like systemd by the calling
func initContainerLabels() (string, string) {
	return "", ""
}

/*
containerLabels returns an allocated processLabel and fileLabel to be used for
container labeling by the calling process.
*/
func containerLabels() (processLabel string, fileLabel string) {
	return "", ""
}

// securityCheckContext validates that the SELinux label is understood by the kernel
func securityCheckContext(val string) error {
	return nil
}

/*
copyLevel returns a label with the MLS/MCS level from src label replaced on
the dest label.
*/
func copyLevel(src, dest string) (string, error) {
	return "", nil
}

// chcon changes the `fpath` file object to the SELinux label `label`.
// If `fpath` is a directory and `recurse`` is true, Chcon will walk the
// directory tree setting the label.
func chcon(fpath string, label string, recurse bool) error {
	return nil
}

// dupSecOpt takes an SELinux process label and returns security options that
// can be used to set the SELinux Type and Level for future container processes.
func dupSecOpt(src string) ([]string, error) {
	return nil, nil
}

// disableSecOpt returns a security opt that can be used to disable SELinux
// labeling support for future container processes.
func disableSecOpt() []string {
	return []string{"disable"}
}
