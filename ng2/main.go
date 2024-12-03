package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jxskiss/mcli"
)

func main() {
	app := &mcli.App{
		Description: "a static festival timetable app builder",
	}
	app.AddRoot(deploy)
	app.Add("bundle", bundle, "only for feti development: rebuild style in static from sass source")
	app.Add("new", new, "create new timetable project in the current folder")
	app.Run()
}

func new() {
	var args struct {
		Name            string `cli:"#R, name, project 'name'"`
		TimetableSource string `cli:"-s, --source, source of timetable data 'csv/json/baserow/pretalx'" default:"json"`
	}
	mcli.Parse(&args)
	if !slices.Contains(DataSources, args.TimetableSource) {
		fmt.Printf("unknown timetable source '%s', use 'json/baserow/pretalx' instead\n", args.TimetableSource)
		os.Exit(1)
	}
	prj, err := ExampleProject(args.Name, strings.ToLower(args.TimetableSource))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = prj.Create()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func bundle() {
	BuildSass(false)
}

func deploy() {
	var args struct {
		ConfigPath string `cli:"#R, -c, --config, path to config file"`
		OutputDir  string `cli:"#R, -o, --output, output directory"`
		LiveServe  bool   `cli:"-s, --serve, serve result with live-rebuild for development"`
	}
	mcli.Parse(&args)
	if args.LiveServe {
		fmt.Println("Please note that this live server is for development purposes only.")
	}

	info, err := os.Stat(args.OutputDir)
	if os.IsNotExist(err) {
		fmt.Printf("given output dir '%s' does not exist\n", args.OutputDir)
		os.Exit(1)
	} else if !info.IsDir() {
		fmt.Printf("given output path '%s' is not a dir\n", args.OutputDir)
		os.Exit(1)
	}

	cfg, err := ConfigFromFile(args.ConfigPath)
	handleErr(err)

	handleErr(cfg.Validate())

	dpl := Deploy{
		Config:    *cfg,
		OutputDir: args.OutputDir,
		LiveServe: args.LiveServe,
		ConfigDir: filepath.Dir(args.ConfigPath),
	}
	handleErr(dpl.Build())

	if args.LiveServe {
		handleErr(dpl.Serve())
	}
}

func handleErr(err error) {
	if err == nil {
		return
	}
	fmt.Println(err)
	os.Exit(1)
}
