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
	AssetsDir        string  `toml:"assets_dir"`
	Genres           []Genre `toml:"genres"`
	TimetableSource  string  `toml:"timetable_source"`
	TimetableJSON    string  `toml:"timetable_json"`
	TimetableCSV     string  `toml:"timetable_csv"`
	BaserowToken     string  `toml:"baserow_token"`
	PretalxToken     string  `toml:"pretalx_token"`
}

func ExampleConfig(timetableSource string) (Config, error) {
	rsl := Config{
		EventName:        "Fetival 2025",
		EventDescription: "The Fetival is like the best festival ever.",
		Logo:             "logo.svg",
		Favicon:          "favicon.svg",
		OpenGraphImage:   "open-graph.png",
		InfoPage:         "infos.md",
		AssetsDir:        "assets",
		Genres:           []Genre{ExampleGenre},
		TimetableSource:  timetableSource,
	}
	if timetableSource == "json" {
		rsl.TimetableJSON = "schedule.json"
	} else if timetableSource == "csv" {
		rsl.TimetableCSV = "schedule.csv"
	} else if timetableSource == "baserow" {
		rsl.BaserowToken = "<MY-SECRET-TOKEN>"
	} else if timetableSource == "pretalx" {
		rsl.PretalxToken = "<MY-SECRET-TOKEN>"
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

func (c Config) Validate() error {
	if c.AssetsDir == "" {
		return fmt.Errorf("field assets_dir is not set in config")
	}

	warn := func(source, field string) {
		fmt.Printf("Warning: %s is used as source for timetable but field '%s' is set in config.\n", source, field)
	}
	errMsg := func(sourceType, field string) error {
		return fmt.Errorf("timetable data source set to '%s' in config file but '%s' field is undefined", sourceType, field)
	}
	switch c.TimetableSource {
	case DataSourceJSON:
		if c.TimetableCSV != "" {
			warn("JSON file", "timetable_csv")
		} else if c.BaserowToken != "" {
			warn("JSON file", "baserow_token")
		} else if c.PretalxToken != "" {
			warn("JSON file", "pretalx_token")
		}
		if c.TimetableJSON == "" {
			return errMsg("json", "timetable_json")
		}
	case DataSourceCSV:
		if c.TimetableJSON != "" {
			warn("CSV file", "timetable_json")
		} else if c.BaserowToken != "" {
			warn("CSV file", "baserow_token")
		} else if c.PretalxToken != "" {
			warn("CSV file", "pretalx_token")
		}
		if c.TimetableCSV == "" {
			return errMsg("csv", "timetable_csv")
		}
	case DataSourceBaserow:
		if c.TimetableJSON != "" {
			warn("Baserow", "timetable_json")
		} else if c.TimetableCSV != "" {
			warn("Baserow", "timetable_csv")
		} else if c.PretalxToken != "" {
			warn("Baserow", "pretalx_token")
		}
		if c.BaserowToken == "" {
			return errMsg("Baserow", "baserow_token")
		}
	case DataSourcePretalx:
		if c.TimetableJSON != "" {
			warn("Pretalx", "timetable_json")
		} else if c.TimetableCSV != "" {
			warn("Pretalx", "timetable_csv")
		} else if c.BaserowToken != "" {
			warn("Pretalx", "baserow_token")
		}
		if c.PretalxToken == "" {
			return errMsg("Pretalx", "pretalx_token")
		}
	}
	return nil
}
