from pathlib import Path
import shutil
import tempfile

from baserow.field import GlobalClient

from feti.config import config, secrets


FAVICON_PATH = Path("assets/img/favicon.svg")
LOGO_PATH = Path("assets/img/logo.svg")
OG_IMAGE_PATH = Path("assets/img/open_graph.png")
MAP_PATH = Path("assets/img/map.png")


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
        self.deploy_folder()
        self.deploy_folder()

    def copy_skeleton(self):
        self.copy_folder_content(self.web_folder_path(), self.path)
        self.replace_file_from_config(
            config().logo,
            self.path / LOGO_PATH,
            "logo",
        )
        self.replace_file_from_config(
            config().favicon,
            self.path / FAVICON_PATH,
            "favicon",
        )
        self.replace_file_from_config(
            config().open_graph_image,
            self.path / OG_IMAGE_PATH,
            "og_image",
        )
        self.replace_file_from_config(
            config().map_source,
            self.path / MAP_PATH,
            "map",
        )

    def replace_file_from_config(self, config_value: str, destination: Path, name: str):
        if config_value == "":
            return
        Deploy.replace_file(
            config().as_absolute_path(config_value),
            destination,
            name,
        )

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
        )
        schedule.sort_schedule()
        with open(self.path / "schedule.json", "w") as f:
            f.write(schedule.model_dump_json())
        await GlobalClient().close()

    def deploy_folder(self):
        self.destination.mkdir(parents=True, exist_ok=True)
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
    def replace_file(source: Path, destination: Path, name: str):
        if source.suffix != destination.suffix:
            raise ValueError(
                f"cannot use {source} for {name} has to be {destination.suffix}",  # noqa
            )
        shutil.copy(source, destination)

    @staticmethod
    def clear_folder(folder: Path):
        for item in folder.iterdir():
            if item.is_dir():
                shutil.rmtree(item)
            else:
                item.unlink()
