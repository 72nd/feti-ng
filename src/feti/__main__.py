import argparse

from livereload import Server


def generate(config: str, secrets: str, output: str):
    """
    Generate the schedule.json based on Baserow data.
    """


def serve():
    """
    Serve page with livereload for development purposes.
    """
    server = Server()
    server.watch("web/**/*")
    server.serve()


def main():
    parser = argparse.ArgumentParser(
        description="takes data from a Baserow instance and generates a schedule.json"
    )
    parser.add_argument(
        "-c",
        "--config",
        required=True,
        help="Path to config file",
    )
    parser.add_argument(
        "-s",
        "--secrets",
        required=True,
        help="Path to the secrets file.",
    )
    parser.add_argument(
        "-o",
        "--output",
        required=True,
        help="Output file or directory.",
    )

    subparsers = parser.add_subparsers(
        dest="command",
        help="Subcommand to run",
    )

    subparsers.required = False  # Make subcommands optional
    parser_serve = subparsers.add_parser(
        "serve",
        help="starts a local webserver with live-reload for development",
    )
    parser_serve.set_defaults(func=serve)

    args = parser.parse_args()
    if hasattr(args, "func"):
        args.func(args)
    else:
        # Implement the main functionality when no subcommand is called
        generate(args.config, args.secrets, args.output)
