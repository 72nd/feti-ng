package main

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed static/*
var staticFiles embed.FS

//go:embed tpl/*
var templateFiles embed.FS

type Deploy struct {
	Config    Config
	OutputDir string
	LiveServe bool
	ConfigDir string
}

func (d Deploy) Build() error {
	// TODO Using tmp folder.
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
	if err := d.deploySchedule(); err != nil {
		return err
	}
	if err := d.deployHTML(); err != nil {
		return err
	}

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
	return CopyDir(srcPath, dstPath)
}

func (d Deploy) deployStatic() error {
	if d.LiveServe {
		return CopyDir("static", filepath.Join(d.OutputDir, "static"))
	}

	dstPath := filepath.Join(d.OutputDir, "static")
	err := os.MkdirAll(dstPath, os.ModePerm)
	if err != nil {
		return err
	}
	return CopyFS(dstPath, staticFiles, "static")
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
	return CopyFile(srcPath, dstPath, true)
}

func (d Deploy) deploySchedule() error {
	dstPath := filepath.Join(d.OutputDir, "schedule.json")
	switch d.Config.TimetableSource {
	case DataSourceJSON:
		path := filepath.Join(d.ConfigDir, d.Config.TimetableJSON)
		rsl, err := ScheduleFromJSON(path, d.Config.Genres)
		if err != nil {
			return err
		}
		if err := rsl.ToJSON(dstPath); err != nil {
			return err
		}
	}
	return nil
}

func (d Deploy) deployHTML() error {
	var tmpl *template.Template
	tmpl = template.New("").Funcs(template.FuncMap{
		"defaultI18nConfig": func() I18nConfig { return d.Config.DefaultI18nConfig() },
	})
	if d.LiveServe {
		var err error
		tmpl, err = tmpl.ParseGlob(filepath.Join("tpl", "*.tmpl.html"))
		if err != nil {
			return err
		}

	} else {
		var err error
		tmpl, err = tmpl.ParseFS(templateFiles, "tpl/*.tmpl.html")
		if err != nil {
			return err
		}
	}

	file, err := os.Create(filepath.Join(d.OutputDir, "index.html"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.ExecuteTemplate(file, "index.tmpl.html", d.Config)
}

func (d Deploy) WatchFiles() {
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
