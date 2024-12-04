package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Like the os.CopyFS() but with the option to set the root.
func CopyFS(dir string, fsSys embed.FS, root string) error {
	return fs.WalkDir(fsSys, root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dir, relPath)

		if entry.IsDir() {
			err := os.MkdirAll(dstPath, os.ModePerm)
			if err != nil {
				return err
			}
			return nil
		}

		data, err := fsSys.ReadFile(path)
		if err != nil {
			return err
		}
		err = os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)
		if err != nil {
			return err
		}
		return os.WriteFile(dstPath, data, os.ModePerm)
	})
}

// Copy a single file.
func CopyFile(srcPath, dstPath string, createParentFolders bool) error {
	if createParentFolders {
		err := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)
		if err != nil {
			return err
		}
	}
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	return os.WriteFile(dstPath, data, os.ModePerm)
}

// Copy a dir including all sub-dirs.
func CopyDir(srcPath, dstPath string) error {
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return err
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("CopyFolder: '%s' is not a dir", srcPath)
	}

	entries, err := os.ReadDir(srcPath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		src := filepath.Join(srcPath, entry.Name())
		dst := filepath.Join(dstPath, entry.Name())

		if entry.IsDir() {
			if err := CopyDir(src, dst); err != nil {
				return err
			}
			continue
		}
		if err := CopyFile(src, dst, false); err != nil {
			return err
		}
	}
	return nil
}
