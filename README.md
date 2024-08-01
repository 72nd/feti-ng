# feti-ng

A client-side HTML5 festival timetable with a static site generator. It uses AlpineJS to render `schedule.json` into a responsive timetable that can be hosted on any simple HTTP server.

A Python script generates the `schedule.json` from data stored in our Baserow instance.

## Usage

Serve the content of the `web` directory using your preferred HTTP server. To modify the schedule or other content, edit the files within the `web` directory.

For more complex deployment needs, use the `feti` Python application, which automates the building and deployment of the web folder, sourcing schedule data from Baserow.

## Site Builder

This Python script generates a web folder from Baserow data and a configuration file, which can then be deployed to a specified location.

### Installation

```bash
python3 -m venv .venv
source .venv/bin/activate
pip3 install .
```

### Baserow Structure

- **Table `Bewerbung`** (entry)
  - `Name`: Text field for artist name
  - `Genre`: Single select field
  - `Titel`: Text field
  - `Dauer`: Duration field
  - `Beschreibung`: Text field
- **Table `Spielstätte`**
  - `Name`: Text field
- **Table `Timetable`**
  - `Wann`: Datetime field
  - `Beitrag`: Link field (to registration)
  - `Spielstätte`: Link field (to location)
  - `Permanent?`: Boolean field (for installations or exhibitions lasting the entire festival)

Field names can be customized using the configuration file.

### Config File

Configuration is done in `config.toml` and `secrets.toml`. Templates for these files are located in the root of the repository. Create a new project folder with copies of these templates and adjust them as needed.

```bash
mkdir my-event
cp config.tpl.toml my-event/config.toml
cp secrets.tpl.toml my-event/secrets.toml
```

You can customize images (logos, etc.) by setting paths in the config file. If paths are left empty, default images will be used. Paths are relative to the config file. Example:

```toml
event_name = "My Event – A Festival"
event_description = "It's like the best thing ever"
logo = "logo.svg"
favicon = "favicon.svg"
open_graph_image = "og.png"

baserow_url = "https://my.baserow.com"
[...]
```

### Usage

```bash
feti -c my-event/config.toml -s my-event/secrets.toml -o /var/www/html
```

## Live-Reload Server for Development

To speed up development, you can use a local web server with live-reload:

```bash
pip3 install .[dev]
feti serve /var/www/data
```