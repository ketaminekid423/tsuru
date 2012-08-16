package fs

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

// FakeFile representss a fake instance of the File interface.
//
// Methods from FakeFile act like methods in os.File, but instead of working in
// a real file, them work in an internal string.
//
// An instance of FakeFile is returned by RecordingFs.Open method.
type FakeFile struct {
	content string
	r       *strings.Reader
}

func (f *FakeFile) reader() *strings.Reader {
	if f.r == nil {
		f.r = strings.NewReader(f.content)
	}
	return f.r
}

func (f *FakeFile) Close() error {
	return nil
}

func (f *FakeFile) Read(p []byte) (n int, err error) {
	return f.reader().Read(p)
}

func (f *FakeFile) ReadAt(p []byte, off int64) (n int, err error) {
	return f.reader().ReadAt(p, off)
}

func (f *FakeFile) Seek(offset int64, whence int) (int64, error) {
	return f.reader().Seek(offset, whence)
}

func (f *FakeFile) Stat() (fi os.FileInfo, err error) {
	return
}

// RecordingFs implements the Fs interface providing a "recording" file system.
//
// A recording file system is a file system that does not execute any action,
// just record them.
//
// All methods from RecordingFs never return errors.
type RecordingFs struct {
	actions []string

	// FileContent is used to provide content for files opened using
	// RecordingFs.
	FileContent string
}

// HasAction checks if a given action was executed in the filesystem.
//
// For example, when you call the Open method with the "/tmp/file.txt"
// argument, RecordingFs will store locally the action "open /tmp/file.txt" and
// you can check it calling HasAction:
//
//     rfs.Open("/tmp/file.txt")
//     rfs.HasAction("open /tmp/file.txt") // true
func (r *RecordingFs) HasAction(action string) bool {
	for _, a := range r.actions {
		if action == a {
			return true
		}
	}
	return false
}

func (r *RecordingFs) Create(name string) (File, error) {
	r.actions = append(r.actions, "create "+name)
	fil := FakeFile{content: r.FileContent}
	return &fil, nil
}

func (r *RecordingFs) Mkdir(name string, perm os.FileMode) error {
	r.actions = append(r.actions, fmt.Sprintf("mkdir %s with mode %#o", name, perm))
	return nil
}

func (r *RecordingFs) MkdirAll(path string, perm os.FileMode) error {
	r.actions = append(r.actions, fmt.Sprintf("mkdirall %s with mode %#o", path, perm))
	return nil
}

// Open returns a FakeFile. The content of the file is provided by the
// FileContent field.
func (r *RecordingFs) Open(name string) (File, error) {
	r.actions = append(r.actions, "open "+name)
	fil := FakeFile{content: r.FileContent}
	return &fil, nil
}

func (r *RecordingFs) Remove(name string) error {
	r.actions = append(r.actions, "remove "+name)
	return nil
}

func (r *RecordingFs) RemoveAll(path string) error {
	r.actions = append(r.actions, "removeall "+path)
	return nil
}

func (r *RecordingFs) Stat(name string) (os.FileInfo, error) {
	r.actions = append(r.actions, "stat "+name)
	return nil, nil
}

// FailureFs is like RecordingFs, except that it returns ENOENT on Open.
type FailureFs struct {
	RecordingFs
}

// Open is used to simulate ENOENT.
func (r *FailureFs) Open(name string) (File, error) {
	r.RecordingFs.Open(name)
	err := os.PathError{
		Err: syscall.ENOENT,
	}
	return nil, &err
}
