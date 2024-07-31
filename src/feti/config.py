from pydantic import BaseModel

import sys

if sys.version_info >= (3, 11):
    import tomllib
else:
    import tomli as tomllib


class Config(BaseModel):
    entry_table_id: int
    entry_name_field: str
    entry_genre_field: str
    entry_title_field: str
    entry_duration_field: str
    entry_description_field: str
    location_table_id: int
    location_name_field: str
    timetable_table_id: int
    timetable_datetime_field: str
    timetable_entry_field: str
    timetable_location_field: str
    timetable_is_permanent_field: str

    @classmethod
    def from_file(cls, path: str):
        with open(path, "rb") as f:
            data = tomllib.load(f)
        return cls(**data)


class Secrets(BaseModel):
    baserow_token: str

    @classmethod
    def from_file(cls, path: str):
        with open(path, "rb") as f:
            data = tomllib.load(f)
        return cls(**data)
