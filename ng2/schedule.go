package main

import (
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

type Schedule struct {
	Genres    []Genre                       `json:"genres"`
	CreatedOn time.Time                     `json:"created_on"`
	Permanent []ScheduleEntry               `json:"permanent"`
	PerDay    map[time.Time][]ScheduleEntry `json:"per_day"`
}
