package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

const (
	DataSourceJSON    = "json"
	DataSourceBaserow = "baserow"
	DataSourcePretalx = "pretalx"
)

var DataSources = []string{
	DataSourceJSON,
	DataSourceBaserow,
	DataSourcePretalx,
}

//go:embed prj/*
var projectFiles embed.FS

type Project struct {
	Name   string
	Config Config
}

func ExampleProject(name, timetableSource string) (Project, error) {
	cfg, err := ExampleConfig(timetableSource)
	if err != nil {
		return Project{}, err
	}
	return Project{
		Name:   name,
		Config: cfg,
	}, nil
}

func (p Project) Create() error {
	if err := p.createFolder(); err != nil {
		return err
	}
	if err := p.createConfig(); err != nil {
		return err
	}
	return fs.WalkDir(projectFiles, "prj", p.populate)
}

func (p Project) createFolder() error {
	path := p.Path()
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("folder '%s' already exists", path)
	} else if !os.IsNotExist(err) {
		return err
	}
	return os.Mkdir(path, os.ModePerm)
}

func (p Project) createConfig() error {
	path := path.Join(p.Path(), "config.toml")
	return p.Config.ToFile(path)
}

func (p Project) populate(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	relPath, err := filepath.Rel("prj", path)
	if err != nil {
		return err
	}
	destPath := filepath.Join(p.Path(), relPath)

	if d.IsDir() {
		return os.MkdirAll(destPath, os.ModePerm)
	}

	data, err := projectFiles.ReadFile(path)
	if err != nil {
		return err
	}
	return os.WriteFile(destPath, data, os.ModePerm)
}

func (p Project) Path() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path.Join(cwd, p.Name)
}
