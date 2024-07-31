from typing import Optional
from livereload.watcher import glob
from pydantic import BaseModel

import sys

if sys.version_info >= (3, 11):
    import tomllib
else:
    import tomli as tomllib


__config: Optional["Config"] = None
__secrets: Optional["Secrets"] = None


class Config(BaseModel):
    baserow_url: str
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


def load_config(path: str):
    global __config
    __config = Config.from_file(path)


def config() -> Config:
    global __config
    if __config is None:
        raise RuntimeError("config was not loaded")
    return __config


class Secrets(BaseModel):
    baserow_token: str

    @classmethod
    def from_file(cls, path: str):
        with open(path, "rb") as f:
            data = tomllib.load(f)
        return cls(**data)


def load_secrets(path: str):
    global __secrets
    __secrets = Secrets.from_file(path)


def secrets() -> Secrets:
    global __secrets
    if __secrets is None:
        raise RuntimeError("secrets was not loaded")
    return __secrets
