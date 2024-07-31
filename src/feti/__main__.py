import argparse
import asyncio
from pathlib import Path
from typing import Optional

from baserow.client import GlobalClient
from baserow.filter import AndFilter
from livereload import Server

from feti.config import config, load_config, load_secrets, secrets


async def generate(
    config_path: str,
    secrets_path: str,
    output: str,
    event_name: Optional[str],
):
    """
    Generate the schedule.json based on Baserow data.
    """
    load_config(config_path)
    load_secrets(secrets_path)
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
        "Event Name" if event_name is not None else "TODO <Event Name>",
    )
    schedule.sort_schedule()
    with open(output, "w") as f:
        f.write(schedule.model_dump_json())
    await GlobalClient().close()


def serve():
    """
    Serve page with livereload for development purposes.
    """
    this_files_path = Path(__file__).resolve()
    root_path = this_files_path.parents[2]
    web_folder = root_path / "web"

    server = Server()
    server.watch(web_folder)
    server.serve(root=web_folder)


def main():
    parser = argparse.ArgumentParser(
        description="takes data from a Baserow instance and generates a schedule.json"
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
    parser.add_argument(
        "--event-name",
        help="Optional name of the event",
    )

    subparsers = parser.add_subparsers(
        dest="command",
        help="Subcommand to run",
    )
    parser_serve = subparsers.add_parser(
        "serve",
        help="starts a local webserver with live-reload for development",
    )
    parser_serve.set_defaults(func=serve)

    args = parser.parse_args()

    if args.command == "serve":
        args.func()
    elif args.config and args.secrets and args.output:
        asyncio.run(generate(
            args.config,
            args.secrets,
            args.output,
            args.event_name,
        ))
    else:
        parser.print_help()
