package fsutil

import (
	"errors"
	"os"
	"path/filepath"
)

type dir_entry struct {
	path string
	next *dir_entry
}

var (
	dir_stack *dir_entry = nil
)

func PushD(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	if err := os.Chdir(abs); err != nil {
		return err
	}

	next := dir_stack
	entry := &dir_entry{
		path: abs,
		next: next,
	}

	dir_stack = entry

	return nil
}

func PopD() error {
	if dir_stack == nil {
		return errors.New("popd: empty dir stack")
	}

	dir_stack = dir_stack.next
	path := dir_stack.path

	return os.Chdir(path)
}
