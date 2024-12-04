package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type ScheduleEntry struct {
	IsPermanent bool          `json:"is_permanent" csv:"is_permanent"`
	StartsAt    time.Time     `json:"starts_at" csv:"starts_at"`
	Duration    time.Duration `json:"duration" csv:"duration"`
	ArtistName  string        `json:"artist_name" csv:"artist_name"`
	Title       string        `json:"title" csv:"title"`
	Description string        `json:"description" csv:"description"`
	Genre       string        `json:"genre" csv:"genre"`
	Location    string        `json:"location" csv:"location"`
}

func (e ScheduleEntry) Validate() error {
	if !e.IsPermanent && e.Duration < 1 {
		return fmt.Errorf("validation error for '%s' entry: duration is not set", e.Title)
	}
	return nil
}

type ScheduleEntries []ScheduleEntry

func (e ScheduleEntries) Validate() error {
	for _, entry := range e {
		if err := entry.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type Schedule struct {
	Genres    []Genre                       `json:"genres"`
	CreatedOn time.Time                     `json:"created_on"`
	Permanent []ScheduleEntry               `json:"permanent"`
	PerDay    map[time.Time][]ScheduleEntry `json:"per_day"`
}

func ScheduleFromJSON(path string, genres []Genre) (*Schedule, error) {
	data, err := ScheduleFileFromJSON(path)
	if err != nil {
		return nil, err
	}
	perDay := data.ScheduledPerDay()
	return &Schedule{
		Genres:    genres,
		CreatedOn: time.Now(),
		PerDay:    perDay,
	}, nil
}

func (s Schedule) ToJSON(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(s)
}
