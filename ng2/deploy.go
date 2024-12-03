package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

// go:embed static/*
// var staticFiles embed.FS

type Deploy struct {
	Config        Config
	OutputDir     string
	LiveServe     bool
	DoRebuildSass bool
}

func (d Deploy) Run() {
	if d.DoRebuildSass {
		d.rebuildSass()
	}
	// Copy static folder
	// Copy assets folder
	// Build html

	// Handle LiveServe 1/2: Rebuild on file change
	// Handle LiveServe 2/2: Start live server
	// --> Maybe it makes sense to split the live control logic to the caller.
}

func (d Deploy) rebuildSass() {

}

func BuildSass(watch bool) error {
	_, err := exec.LookPath("sass")
	if err != nil {
		return fmt.Errorf("sass command not found. Please ensure sass is installed and in your PATH.")
	}
	cmd := exec.Command(
		"sass",
		"sass/bootstrap.scss",
		"static/css/bootstrap.min.css",
		"--style=compressed",
		"--quiet-deps",
	)
	if watch {
		cmd.Args = append(cmd.Args, "--watch")
	}
	return cmd.Run()
}

func ExpandEmbedFS(embedFS embed.FS, fsName, dest string) error {
	return fs.WalkDir(embedFS, fsName, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel("prj", path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dest, relPath)

		if d.IsDir() {
			return os.MkdirAll(destPath, os.ModePerm)
		}

		data, err := projectFiles.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(destPath, data, os.ModePerm)
	})
}
