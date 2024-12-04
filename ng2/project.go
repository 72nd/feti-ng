package main

import (
	"embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

const (
	DataSourceJSON    = "json"
	DataSourceCSV     = "csv"
	DataSourceBaserow = "baserow"
	DataSourcePretalx = "pretalx"
)

var DataSources = []string{
	DataSourceJSON,
	DataSourceCSV,
	DataSourceBaserow,
	DataSourcePretalx,
}

//go:embed prj/*
var projectFiles embed.FS

type Project struct {
	Name            string
	TimetableSource string
	Config          Config
}

func ExampleProject(name, timetableSource string) (Project, error) {
	cfg, err := ExampleConfig(timetableSource)
	if err != nil {
		return Project{}, err
	}
	return Project{
		Name:            name,
		TimetableSource: timetableSource,
		Config:          cfg,
	}, nil
}

func (p Project) Create() error {
	if err := p.createFolder(); err != nil {
		return err
	}
	if err := p.createConfig(); err != nil {
		return err
	}
	// TODO: HIER IST DAS PROBLEM
	if err := os.CopyFS(p.Path(), projectFiles); err != nil {
		return err
	}
	return p.createScheduleSource()
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

func (p Project) createScheduleSource() error {
	switch p.TimetableSource {
	case DataSourceJSON:
		return ExampleScheduleFile.ToJSON(filepath.Join(p.Path(), "schedule.json"))
	case DataSourceCSV:
		return ExampleScheduleFile.ToCSV(filepath.Join(p.Path(), "schedule.csv"))
	}
	return nil
}

func (p Project) Path() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(cwd, p.Name)
}
