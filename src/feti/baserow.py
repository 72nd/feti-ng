from datetime import datetime, timedelta
from typing import Any, Optional

from baserow.field import SingleSelectField
from baserow.table import Table, TableLinkField
from pydantic import ConfigDict, Field

from feti.config import config


class Entry(Table):
    table_id = config().entry_table_id
    table_name = "Entry"
    model_config = ConfigDict(populate_by_name=True)

    name: Optional[str] = Field(
        default=None,
        alias=config().entry_name_field,
    )
    genre: Optional[SingleSelectField[Any]] = Field(
        default=None,
        alias=config().entry_genre_field,
    )
    title: Optional[str] = Field(
        default=None,
        alias=config().entry_title_field,
    )
    duration: Optional[timedelta] = Field(
        default=None,
        alias=config().entry_duration_field,
    )
    description: Optional[str] = Field(
        default=None,
        alias=config().entry_description_field,
    )


class Location(Table):
    table_id = config().location_table_id
    table_name = "Location"
    model_config = ConfigDict(populate_by_name=True)

    name: str = Field(alias=config().location_name_field)


class Timetable(Table):
    table_id = config().timetable_table_id
    table_name = "Timetable"
    model_config = ConfigDict(populate_by_name=True)

    starts_at: Optional[datetime] = Field(
        default=None,
        alias=config().timetable_datetime_field,
    )
    entry: TableLinkField[Entry] = Field(
        alias=config().timetable_entry_field,
    )
    location: TableLinkField[Location] = Field(
        alias=config().timetable_location_field,
    )
    is_permanent: bool = Field(alias=config().timetable_is_permanent_field)

    def debug_str(self) -> str:
        return f"timetable entry with id {self.row_id} at {self.starts_at}"
