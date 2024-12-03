package main

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed static/*
var staticFiles embed.FS

type Deploy struct {
	Config    Config
	OutputDir string
	LiveServe bool
	ConfigDir string
}

func (d Deploy) Build() error {
	if d.LiveServe {
		BuildSass(false)
	}
	if err := d.deployAssets(); err != nil {
		return err
	}
	if err := d.deployStatic(); err != nil {
		return err
	}
	if err := d.deploySysAssets(); err != nil {
		return err
	}

	// Copy static folder
	// Copy assets folder
	// Build html

	// Handle LiveServe 1/2: Rebuild on file change
	// Handle LiveServe 2/2: Start live server
	// --> Maybe it makes sense to split the live control logic to the caller.
	return nil
}

func (d Deploy) deployAssets() error {
	srcPath := filepath.Join(d.ConfigDir, d.Config.AssetsDir)
	dstPath := filepath.Join(d.OutputDir, "assets")
	err := os.MkdirAll(dstPath, os.ModePerm)
	if err != nil {
		return err
	}
	return os.CopyFS(dstPath, os.DirFS(srcPath))
}

func (d Deploy) deployStatic() error {
	if d.LiveServe {
		return os.CopyFS(d.OutputDir, os.DirFS("static"))
	}

	dstPath := filepath.Join(d.OutputDir, "static")
	err := os.MkdirAll(dstPath, os.ModePerm)
	if err != nil {
		return err
	}
	return os.CopyFS(dstPath, staticFiles)
}

func (d Deploy) deploySysAssets() error {
	items := [][]string{
		{d.Config.Favicon, "static/img/favicon.svg"},
		{d.Config.Logo, "static/img/logo.svg"},
		{d.Config.OpenGraphImage, "static/img/open-graph.png"},
	}
	for _, item := range items {
		err := d.copySysAsset(item[0], item[1])
		if err != nil {
			return err
		}
	}
	return nil
}

func (d Deploy) copySysAsset(configPathValue, dstPath string) error {
	srcPath := filepath.Join(d.ConfigDir, configPathValue)
	dstPath = filepath.Join(d.OutputDir, dstPath)
	err := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	return os.WriteFile(dstPath, data, os.ModePerm)
}

func (d Deploy) WatchFiles() {
}

func (d Deploy) Serve() error {
	return nil
}

func BuildSass(watch bool) error {
	_, err := exec.LookPath("sass")
	if err != nil {
		return fmt.Errorf("sass command not found. Please ensure sass is installed and in your PATH")
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
