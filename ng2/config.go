package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Genre struct {
	Name  string `toml:"name"`
	Color string `toml:"color"`
	Icon  string `toml:"icon"`
}

var ExampleGenre = Genre{
	Name:  "Music",
	Color: "genre-color-vermilion",
	Icon:  "music_note",
}

type Config struct {
	EventName        string  `toml:"event_name"`
	EventDescription string  `toml:"event_description"`
	Logo             string  `toml:"logo"`
	Favicon          string  `toml:"favicon"`
	OpenGraphImage   string  `toml:"open_graph_image"`
	InfoPage         string  `toml:"info_page"`
	AssetsFolder     string  `toml:"assets_folder"`
	Genres           []Genre `toml:"genres"`
	TimetableSource  string  `toml:"timetable_source"`
	TimetableJSON    string  `toml:"timetable_json"`
	BaserowToken     string  `toml:"baserow_token"`
}

func ExampleConfig(timetableSource string) (Config, error) {
	rsl := Config{
		EventName:        "Fnordival 2025",
		EventDescription: "The Fnordival is like the best festival ever.",
		Logo:             "logo.svg",
		Favicon:          "favicon.svg",
		OpenGraphImage:   "open_graph.png",
		InfoPage:         "infos.md",
		AssetsFolder:     "assets",
		Genres:           []Genre{ExampleGenre},
		TimetableSource:  timetableSource,
	}
	if timetableSource == "json" {
		rsl.TimetableJSON = "schedule.json"
	} else if timetableSource == "baserow" {
		rsl.BaserowToken = "<MY-SECRET-TOKEN>"
	} else {
		return Config{}, fmt.Errorf("timetable source '%s' is currently not implemented", timetableSource)
	}
	return rsl, nil
}

func ConfigFromFile(path string) (*Config, error) {
	var rsl Config
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	_, err = toml.Decode(string(file), &rsl)
	return &rsl, err
}

func (c Config) ToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	err = toml.NewEncoder(f).Encode(c)
	if err != nil {
		return err
	}
	return f.Close()
}
