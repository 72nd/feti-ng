package main

import (
	"encoding/json"
	"os"
	"sort"
	"time"

	"github.com/gocarina/gocsv"
)

// Use a JSON or CSV file as data source.
type ScheduleFile ScheduleEntries

var ExampleScheduleFile = ScheduleFile{
	{
		IsPermanent: true,
		StartsAt:    time.Time{},
		Duration:    0,
		ArtistName:  "Festival Info Booth",
		Title:       "Information Booth",
		Description: "Your go-to spot for festival maps, schedules, and assistance.",
		Genre:       "Info",
		Location:    "Main Entrance",
	},
	{
		IsPermanent: false,
		StartsAt:    time.Date(2024, 12, 3, 18, 0, 0, 0, time.UTC),
		Duration:    2 * time.Hour,
		ArtistName:  "The Groove Masters",
		Title:       "Sunset Groove",
		Description: "An energetic performance by The Groove Masters to kick off the evening.",
		Genre:       "Music",
		Location:    "Main Stage",
	},
	{
		IsPermanent: false,
		StartsAt:    time.Date(2024, 12, 4, 20, 30, 0, 0, time.UTC),
		Duration:    90 * time.Minute,
		ArtistName:  "Lunar Ensemble",
		Title:       "Moonlit Melodies",
		Description: "A serene musical experience under the stars.",
		Genre:       "Music",
		Location:    "Outdoor Theater",
	},
}

func ScheduleFileFromJSON(path string) (ScheduleFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rsl ScheduleFile
	err = json.Unmarshal(data, &rsl)
	if err != nil {
		return nil, err
	}
	return rsl, ScheduleEntries(rsl).Validate()
}

func (s ScheduleFile) ToJSON(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(s)
}

func (s ScheduleFile) ToCSV(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gocsv.MarshalFile(s, f)
}

func (s ScheduleFile) ScheduledPerDay() map[time.Time][]ScheduleEntry {
	rsl := make(map[time.Time][]ScheduleEntry)
	for _, entry := range s {
		if entry.IsPermanent {
			continue
		}
		day := entry.StartsAt.UTC().Truncate(24 * time.Hour)
		rsl[day] = append(rsl[day], entry)
	}

	for day, entries := range rsl {
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].StartsAt.Before(entries[j].StartsAt)
		})
		rsl[day] = entries
	}
	return rsl
}
