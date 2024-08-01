from pathlib import Path
import shutil
import tempfile

from baserow.field import GlobalClient

from feti.config import config, secrets


class Deploy:
    def __init__(self, destination: Path):
        self.dir = tempfile.TemporaryDirectory()
        self.path = Path(self.dir.name)
        self.destination = destination

    def __del__(self):
        self.dir.cleanup()

    async def run(self):
        self.copy_skeleton()
        await self.generate_schedule()

    def copy_skeleton(self):
        self.copy_folder_content(self.web_folder_path(), self.path)

        self.clear_folder(self.destination)
        self.copy_folder_content(self.path, self.destination)

    @staticmethod
    def web_folder_path() -> Path:
        this_file_path = Path(__file__).resolve()
        root_path = this_file_path.parents[2]
        return root_path / "web"

    @staticmethod
    def copy_folder_content(source: Path, destination: Path):
        for item in source.iterdir():
            if item.is_dir():
                shutil.copytree(item, destination / item.name)
            else:
                shutil.copy2(item, destination / item.name)

    @staticmethod
    def clear_folder(folder: Path):
        for item in folder.iterdir():
            if item.is_dir():
                shutil.rmtree(item)
            else:
                item.unlink()

    async def generate_schedule(self):
        GlobalClient.configure(
            config().baserow_url,
            token=secrets().baserow_token,
        )
        from feti.baserow import Entry, Location, Timetable
        from feti.schedule import Schedule
        schedule = Schedule.from_baserow(
            await Entry.query(size=-1),
            await Location.query(size=-1),
            await Timetable.query(size=-1),
            config().event_name,
        )
        schedule.sort_schedule()
        with open(self.path / "schedule.json", "w") as f:
            f.write(schedule.model_dump_json())
        await GlobalClient().close()
