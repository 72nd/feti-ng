from datetime import date, datetime, timedelta

from typing import Optional
from pydantic import BaseModel

from feti.baserow import Entry, Location, Timetable
from feti.config import config


GENRE_COLOR = {
    "Ausstellung": ("genre-color-orange", "wall_art"),
    "Film": ("genre-color-chartreuse", "theaters"),
    "Installation": ("genre-color-yellow", "highlight"),
    "Musik": ("genre-color-vermilion", "music_note"),
    "Performance": ("genre-color-magenta", "stars"),
    "Tanz": ("genre-color-magenta", "directions_walk"),
    "Theater": ("genre-color-teal", "theater_comedy"),
    "Lesung": ("genre-color-blue", "book_4"),
    "Digital": ("genre-color-red", "wifi"),
    "Hybrid": ("genre-color-orange-light", "transition_fade"),
    "Workshop": ("genre-color-green", "handyman"),
    "Kindertheater": ("genre-color-purple", "child_care"),
    "Poetry Slam": ("genre-color-blue-light", "mic_external_on"),
    "DJ": ("genre-color-amber", "nightlife"),
    "Meta": ("genre-color-vermilion-light", "eco"),
}


class ScheduleEntry(BaseModel):
    is_permanent: bool
    starts_at: Optional[datetime] = None
    artist_name: str
    title: str
    duration: Optional[timedelta]
    description: str
    genre: str
    location: str

    @classmethod
    def from_baserow(
        cls,
        timetable_entry: Timetable,
        locations: dict[int, Location],
        entries: dict[int, Entry],
    ):
        entry = cls.entry_for_timetable_entry(timetable_entry, entries)
        if entry.name is None:
            raise ValueError(
                f"{timetable_entry.debug_str()} name field is None",
            )
        if entry.title is None:
            raise ValueError(
                f"{timetable_entry.debug_str()} title field is None",
            )
        if entry.description is None:
            raise ValueError(
                f"{timetable_entry.debug_str()} description field is None",
            )
        if entry.genre is None:
            raise ValueError(
                f"{timetable_entry.debug_str()} genre field is None",
            )
        if entry.genre.value is None:
            raise ValueError(
                f"{timetable_entry.debug_str()} genre.value field is None",
            )

        location = cls.location_for_timetable_entry(timetable_entry, locations)

        return ScheduleEntry(
            is_permanent=timetable_entry.is_permanent,
            starts_at=timetable_entry.starts_at,
            artist_name=entry.name,
            title=entry.title,
            duration=entry.duration,
            description=entry.description,
            genre=entry.genre.value,
            location=location.name,
        )

    @staticmethod
    def entry_for_timetable_entry(timetable_entry: Timetable, entries: dict[int, Entry]) -> Entry:
        if len(timetable_entry.entry.root) > 1:
            raise ValueError(f"{timetable_entry.debug_str()} has more than one linked entries")  # noqa
        elif len(timetable_entry.entry.root) < 1:
            raise ValueError(f"{timetable_entry.debug_str()} has no linked entry")  # noqa
        entry_id = timetable_entry.entry.root[0].row_id
        if not entry_id in entries:
            raise ValueError(f"{timetable_entry.debug_str()} has invalid linked entry with id {entry_id}")  # noqa
        return entries[entry_id]

    @staticmethod
    def location_for_timetable_entry(timetable_entry: Timetable, locations: dict[int, Location]) -> Location:
        if len(timetable_entry.location.root) > 1:
            raise ValueError(f"{timetable_entry.debug_str()} has more than one linked locations")  # noqa
        elif len(timetable_entry.location.root) < 1:
            raise ValueError(f"{timetable_entry.debug_str()} has no linked location")  # noqa
        location_id = timetable_entry.location.root[0].row_id
        if not location_id in locations:
            raise ValueError(f"{timetable_entry.debug_str()} has invalid linked location with id {location_id}")  # noqa
        return locations[location_id]


class Map(BaseModel):
    enabled: bool
    center_y: int
    center_x: int
    zoom: float
    upper_bound_y: int
    upper_bound_x: int

    @classmethod
    def from_config(cls):
        return cls(
            enabled=config().map_enabled,
            center_y=config().map_center_y,
            center_x=config().map_center_x,
            zoom=config().map_zoom,
            upper_bound_y=config().map_upper_bound_y,
            upper_bound_x=config().map_upper_bound_x,
        )


class Schedule(BaseModel):
    event_name: str
    event_description: str
    genre_color_class: dict[str, tuple[str, str]] = {}
    map: Map
    created_on: datetime
    permanent: list[ScheduleEntry] = []
    per_day: dict[date, list[ScheduleEntry]] = {}

    @classmethod
    def from_baserow(
        cls,
        entries: list[Entry],
        locations: list[Location],
        timetable: list[Timetable],
    ):
        entry_dict = {
            item.row_id: item for item in entries if item.row_id is not None
        }
        location_dict = {
            item.row_id: item for item in locations if item.row_id is not None
        }

        rsl = cls(
            event_name=config().event_name,
            event_description=config().event_description,
            genre_color_class=GENRE_COLOR,
            created_on=datetime.now(),
            map=Map.from_config(),
        )

        for tt_entry in timetable:
            tmp = ScheduleEntry.from_baserow(
                tt_entry, location_dict, entry_dict
            )
            if tmp.is_permanent:
                rsl.permanent.append(tmp)
                continue
            else:
                if tmp.starts_at is None:
                    print(
                        f"WARN {tt_entry.debug_str()} has no start datetime set but is not permanent"  # noqa
                    )
                    continue
                day = tmp.starts_at.date()
                if day not in rsl.per_day:
                    rsl.per_day[day] = []
                rsl.per_day[day].append(tmp)
        return rsl

    def sort_schedule(self):
        self.permanent = sorted(self.permanent, key=lambda entry: entry.genre)
        self.per_day = dict(sorted(self.per_day.items()))
        for day, entries in self.per_day.items():
            self.per_day[day] = sorted(
                entries, key=lambda entry: entry.starts_at or datetime.min
            )
