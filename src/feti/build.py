import subprocess
from pathlib import Path
import shutil
import sys


class BuildAssets:
    """
    Abstraction for the building of the assets. This is only needed during
    development.
    """

    def __init__(self, watch: bool):
        base_dir = Path(__file__).resolve().parents[2]
        self.watch = watch
        self.sass_dir = base_dir / "sass"
        self.assets_dir = base_dir / "web/assets"

    def run(self):
        self.__compile_sass()

    def __compile_sass(self):
        self.check_command("sass")
        input_file = self.sass_dir / "bootstrap.scss"
        output_file = self.assets_dir / "css/bootstrap.min.css"
        cmd = [
            "sass",
            str(input_file),
            str(output_file),
            "--style=compressed",
            "--quiet-deps",
        ]
        if self.watch:
            cmd.append("--watch")
        try:
            subprocess.run(cmd, check=True)
        except subprocess.CalledProcessError as e:
            print(f"Failed to compile SCSS: {e}")
            sys.exit(1)

    @staticmethod
    def check_command(cmd: str):
        if not shutil.which(cmd):
            raise EnvironmentError(
                f"{cmd} command not found. Please ensure {cmd} is installed and in your PATH.",  # noqa
            )
