package main

import (
	"fmt"
	"os"
	"path"
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

type Project struct {
	Config Config
}

func ExampleProject(timetableSource string) (Project, error) {
	cfg, err := ExampleConfig(timetableSource)
	if err != nil {
		return Project{}, err
	}
	return Project{
		Config: cfg,
	}, nil
}

func (p Project) Create(name string) error {
	err := p.createFolder(name)
	if err != nil {
		return err
	}
	return nil
}

func (p Project) createFolder(name string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	path := path.Join(cwd, name)
	if _, err = os.Stat(path); err == nil {
		return fmt.Errorf("folder '%s' already exists", path)
	} else if !os.IsNotExist(err) {
		return err
	}
	return os.Mkdir(path, os.ModePerm)
}
