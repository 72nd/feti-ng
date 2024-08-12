import argparse
import asyncio
from pathlib import Path

try:
    from livereload import Server
except ImportError:
    Server = None

from feti.config import load_config, load_secrets
from feti.deploy import Deploy


async def deploy(
    config_path: str,
    secrets_path: str,
    output_path: str,
):
    """Build and deploy the 'static' web content."""
    load_config(config_path)
    load_secrets(secrets_path)
    deploy = Deploy(Path(output_path))
    await deploy.run()


def build(watch: bool):
    """Builds the assets."""


def serve(dir: str):
    """
    Serve page with livereload for development purposes.
    """
    this_files_path = Path(__file__).resolve()
    root_path = this_files_path.parents[2]
    web_folder = root_path / dir
    print(web_folder)

    if Server is None:
        print(
            "additional dependencies needed install with pip3 install feti[dev]")
        return
    server = Server()
    server.watch(web_folder)
    server.serve(root=web_folder)


def main():
    parser = argparse.ArgumentParser(
        description="takes data from a Baserow instance and deploys the web folder"
    )
    parser.add_argument(
        "-c",
        "--config",
        help="Path to config file",
    )
    parser.add_argument(
        "-s",
        "--secrets",
        help="Path to the secrets file.",
    )
    parser.add_argument(
        "-o",
        "--output",
        help="Output file or directory.",
    )

    subparsers = parser.add_subparsers(
        dest="command",
        help="Subcommand to run",
    )

    parser_build = subparsers.add_parser(
        "build",
        help="builds web assets",
    )
    parser_build.add_argument(
        "--watch",
        action="store_true",
        help="enable watch mode to rebuild on file changes",
    )
    parser_build.set_defaults(func=build)

    parser_serve = subparsers.add_parser(
        "serve",
        help="starts a local webserver with live-reload for development",
    )
    parser_serve.add_argument(
        "WEB_DIR",
        help="folder containing the web content"
    )
    parser_serve.set_defaults(func=serve)

    args = parser.parse_args()

    if args.command == "build":
        args.func(args.watch)
    elif args.command == "serve":
        args.func(args.WEB_DIR)
    elif args.config and args.secrets and args.output:
        asyncio.run(deploy(
            args.config,
            args.secrets,
            args.output,
        ))
    else:
        parser.print_help()
