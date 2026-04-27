package selinux

import (
	"github.com/opencontainers/selinux/internal/impl"
)

const (
	// Enforcing constant indicate SELinux is in enforcing mode
	Enforcing = impl.Enforcing
	// Permissive constant to indicate SELinux is in permissive mode
	Permissive = impl.Permissive
	// Disabled constant to indicate SELinux is disabled
	Disabled = impl.Disabled

	// DefaultCategoryRange is the default upper bound on the category range.
	// See [SetCategoryRange].
	DefaultCategoryRange = impl.DefaultCategoryRange
)

var (
	// ErrMCSAlreadyExists is returned when trying to allocate a duplicate MCS.
	ErrMCSAlreadyExists = impl.ErrMCSAlreadyExists
	// ErrEmptyPath is returned when an empty path has been specified.
	ErrEmptyPath = impl.ErrEmptyPath

	// ErrInvalidLabel is returned when an invalid label is specified.
	ErrInvalidLabel = impl.ErrInvalidLabel

	// ErrIncomparable is returned two levels are not comparable
	ErrIncomparable = impl.ErrIncomparable
	// ErrLevelSyntax is returned when a sensitivity or category do not have correct syntax in a level
	ErrLevelSyntax = impl.ErrLevelSyntax

	// ErrContextMissing is returned if a requested context is not found in a file.
	ErrContextMissing = impl.ErrContextMissing
	// ErrVerifierNil is returned when a context verifier function is nil.
	ErrVerifierNil = impl.ErrVerifierNil

	// ErrNotTGLeader is returned by [SetKeyLabel] if the calling thread
	// is not the thread group leader.
	ErrNotTGLeader = impl.ErrNotTGLeader
)

// Context is a representation of the SELinux label broken into 4 parts
type Context = impl.Context

// SetDisabled disables SELinux support for the package
func SetDisabled() {
	impl.SetDisabled()
}

// GetEnabled returns whether SELinux is currently enabled.
func GetEnabled() bool {
	return impl.GetEnabled()
}

// SetCategoryRange allows to adjust the upper bound of the category range.
// It affects subsequent calls to [KVMContainerLabel] and [InitContainerLabel].
func SetCategoryRange(upper uint32) error {
	return impl.SetCategoryRange(upper)
}

// ClassIndex returns the int index for an object class in the loaded policy,
// or -1 and an error
func ClassIndex(class string) (int, error) {
	return impl.ClassIndex(class)
}

// SetFileLabel sets the SELinux label for this path, following symlinks,
// or returns an error.
func SetFileLabel(fpath string, label string) error {
	return impl.SetFileLabel(fpath, label)
}

// LsetFileLabel sets the SELinux label for this path, not following symlinks,
// or returns an error.
func LsetFileLabel(fpath string, label string) error {
	return impl.LsetFileLabel(fpath, label)
}

// FileLabel returns the SELinux label for this path, following symlinks,
// or returns an error.
func FileLabel(fpath string) (string, error) {
	return impl.FileLabel(fpath)
}

// LfileLabel returns the SELinux label for this path, not following symlinks,
// or returns an error.
func LfileLabel(fpath string) (string, error) {
	return impl.LfileLabel(fpath)
}

// SetFSCreateLabel tells the kernel what label to use for all file system objects
// created by this task.
// Set the label to an empty string to return to the default label. Calls to SetFSCreateLabel
// should be wrapped in runtime.LockOSThread()/runtime.UnlockOSThread() until file system
// objects created by this task are finished to guarantee another goroutine does not migrate
// to the current thread before execution is complete.
func SetFSCreateLabel(label string) error {
	return impl.SetFSCreateLabel(label)
}

// FSCreateLabel returns the default label the kernel which the kernel is using
// for file system objects created by this task. "" indicates default.
func FSCreateLabel() (string, error) {
	return impl.ReadConThreadSelf("attr/fscreate")
}

// CurrentLabel returns the SELinux label of the current process thread, or an error.
func CurrentLabel() (string, error) {
	return impl.ReadConThreadSelf("attr/current")
}

// PidLabel returns the SELinux label of the given pid, or an error.
func PidLabel(pid int) (string, error) {
	return impl.PidLabel(pid)
}

// ExecLabel returns the SELinux label that the kernel will use for any programs
// that are executed by the current process thread, or an error.
func ExecLabel() (string, error) {
	return impl.ReadConThreadSelf("attr/exec")
}

// CanonicalizeContext takes a context string and writes it to the kernel
// the function then returns the context that the kernel will use. Use this
// function to check if two contexts are equivalent
func CanonicalizeContext(val string) (string, error) {
	return impl.CanonicalizeContext(val)
}

// ComputeCreateContext requests the type transition from source to target for
// class from the kernel.
func ComputeCreateContext(source string, target string, class string) (string, error) {
	return impl.ComputeCreateContext(source, target, class)
}

// CalculateGlbLub computes the glb (greatest lower bound) and lub (least upper bound)
// of a source and target range.
// The glblub is calculated as the greater of the low sensitivities and
// the lower of the high sensitivities and the and of each category bitset.
func CalculateGlbLub(sourceRange, targetRange string) (string, error) {
	return impl.CalculateGlbLub(sourceRange, targetRange)
}

// SetExecLabel sets the SELinux label that the kernel will use for any programs
// that are executed by the current process thread, or an error. Calls to SetExecLabel
// should  be wrapped in runtime.LockOSThread()/runtime.UnlockOSThread() until execution
// of the program is finished to guarantee another goroutine does not migrate to the current
// thread before execution is complete.
func SetExecLabel(label string) error {
	return impl.WriteConThreadSelf("attr/exec", label)
}

// SetTaskLabel sets the SELinux label for the current thread, or an error.
// This requires the dyntransition permission. Calls to SetTaskLabel should
// be wrapped in runtime.LockOSThread()/runtime.UnlockOSThread() to guarantee
// the current thread does not run in a new mislabeled thread.
func SetTaskLabel(label string) error {
	return impl.WriteConThreadSelf("attr/current", label)
}

// SetSocketLabel takes a process label and tells the kernel to assign the
// label to the next socket that gets created. Calls to SetSocketLabel
// should be wrapped in runtime.LockOSThread()/runtime.UnlockOSThread() until
// the socket is created to guarantee another goroutine does not migrate
// to the current thread before execution is complete.
func SetSocketLabel(label string) error {
	return impl.WriteConThreadSelf("attr/sockcreate", label)
}

// SocketLabel retrieves the current socket label setting
func SocketLabel() (string, error) {
	return impl.ReadConThreadSelf("attr/sockcreate")
}

// PeerLabel retrieves the label of the client on the other side of a socket
func PeerLabel(fd uintptr) (string, error) {
	return impl.PeerLabel(fd)
}

// SetKeyLabel takes a process label and tells the kernel to assign the
// label to the next kernel keyring that gets created.
//
// Calls to SetKeyLabel should be wrapped in
// runtime.LockOSThread()/runtime.UnlockOSThread() until the kernel keyring is
// created to guarantee another goroutine does not migrate to the current
// thread before execution is complete.
//
// Only the thread group leader can set key label.
func SetKeyLabel(label string) error {
	return impl.SetKeyLabel(label)
}

// KeyLabel retrieves the current kernel keyring label setting
func KeyLabel() (string, error) {
	return impl.KeyLabel()
}

// NewContext creates a new Context struct from the specified label
func NewContext(label string) (Context, error) {
	return impl.NewContext(label)
}

// ClearLabels clears all reserved labels
func ClearLabels() {
	impl.ClearLabels()
}

// ReserveLabel reserves the MLS/MCS level component of the specified label.
//
// Deprecated: use [ReserveLabelV2] instead.
func ReserveLabel(label string) {
	_ = impl.ReserveLabel(label)
}

// ReserveLabelV2 reserves the MLS/MCS level component of the specified label.
// Returns an error if the label can't be reserved.
func ReserveLabelV2(label string) error {
	return impl.ReserveLabel(label)
}

// CheckLabel check the MLS/MCS level component of the specified label
func CheckLabel(label string) error {
	return impl.CheckLabel(label)
}

// MLSEnabled checks if MLS is enabled.
func MLSEnabled() bool {
	return impl.MLSEnabled()
}

// EnforceMode returns the current SELinux mode (one of [Enforcing],
// [Permissive], or [Disabled]).
func EnforceMode() int {
	return impl.EnforceMode()
}

// SetEnforceMode sets the current SELinux mode Enforcing, Permissive.
// Disabled is not valid, since this needs to be set at boot time.
func SetEnforceMode(mode int) error {
	return impl.SetEnforceMode(mode)
}

// DefaultEnforceMode returns the systems default SELinux mode Enforcing,
// Permissive or Disabled. Note this is just the default at boot time.
// EnforceMode tells you the systems current mode.
func DefaultEnforceMode() int {
	return impl.DefaultEnforceMode()
}

// ReleaseLabel un-reserves the MLS/MCS Level field of the specified label,
// allowing it to be used by another process.
func ReleaseLabel(label string) {
	impl.ReleaseLabel(label)
}

// ROFileLabel returns the specified SELinux readonly file label.
//
// Deprecated: this (apparently) has no users and will be removed from the
// future version of this package. Open a bug report if you use it.
func ROFileLabel() string {
	return impl.ROFileLabel()
}

// KVMContainerLabels returns the default processLabel and mountLabel to be used
// for kvm containers by the calling process.
//
// Deprecated: use [KVMContainerLabel] instead.
func KVMContainerLabels() (string, string) {
	return impl.KVMContainerLabels()
}

// KVMContainerLabel returns the default process label to be used
// for KVM containers by the calling process.
func KVMContainerLabel() (string, error) {
	return impl.KVMContainerLabel()
}

// InitContainerLabels returns the default processLabel and file labels to be
// used for containers running an init system like systemd by the calling process.
//
// Deprecated: use [InitContainerLabel] instead.
func InitContainerLabels() (string, string) {
	return impl.InitContainerLabels()
}

// InitContainerLabel returns the default process label to be used
// for containers running an init system like systemd by the calling process.
func InitContainerLabel() (string, error) {
	return impl.InitContainerLabel()
}

// ContainerLabels returns an allocated processLabel and fileLabel to be used for
// container labeling by the calling process.
//
// Deprecated: this (apparently) has no users and will be removed from the
// future version of this package. Open a bug report if you use it.
func ContainerLabels() (processLabel string, fileLabel string) {
	return impl.ContainerLabels()
}

// SecurityCheckContext validates that the SELinux label is understood by the kernel
func SecurityCheckContext(val string) error {
	return impl.SecurityCheckContext(val)
}

// CopyLevel returns a label with the MLS/MCS level from src label replaced on
// the dest label.
func CopyLevel(src, dest string) (string, error) {
	return impl.CopyLevel(src, dest)
}

// Chcon changes the fpath file object to the SELinux label.
// If fpath is a directory and recurse is true, then Chcon walks the
// directory tree setting the label.
//
// The fpath itself is guaranteed to be relabeled last.
func Chcon(fpath string, label string, recurse bool) error {
	return impl.Chcon(fpath, label, recurse)
}

// DupSecOpt takes an SELinux process label and returns security options that
// can be used to set the SELinux Type and Level for future container processes.
func DupSecOpt(src string) ([]string, error) {
	return impl.DupSecOpt(src)
}

// DisableSecOpt returns a security opt that can be used to disable SELinux
// labeling support for future container processes.
func DisableSecOpt() []string {
	return []string{"disable"}
}

// SEUserByName retrieves the SELinux username and security level for a given
// Linux username. The username and security level is based on the
// /etc/selinux/{SELINUXTYPE}/seusers file.
func SEUserByName(username string) (seUser string, level string, err error) {
	return impl.SEUserByName(username)
}

// GetDefaultContextWithLevel gets a single context for the specified SELinux user
// identity that is reachable from the specified scon context. The context is based
// on the per-user /etc/selinux/{SELINUXTYPE}/contexts/users/<username> if it exists,
// and falls back to the global /etc/selinux/{SELINUXTYPE}/contexts/default_contexts
// file and finally the global /etc/selinux/{SELINUXTYPE}/contexts/failsafe_context
// file if no match can be found anywhere else.
func GetDefaultContextWithLevel(user, level, scon string) (string, error) {
	return impl.GetDefaultContextWithLevel(user, level, scon)
}

// PrivContainerMountLabel returns mount label for privileged containers
func PrivContainerMountLabel() string {
	return impl.PrivContainerMountLabel()
}
