# feti-ng

HTML5 client side festival timetable with static site generator. Uses AlpineJS to render the `schedule.json` as a responsive timetable. Can be deployed using any simple HTTP server.

The python script is used to generate a `schedule.json` from our Baserow instance.


## Usage

Use your favorite HTTP-server to serve the content of the `web` directory. Done. You can alter the schedule, etc. by editing the files within `web`.

If a more complex deployment is needed you can use the `feti` python application for automatic build and deployment of the web folder. The script uses Baserow as the source for the schedule data.


## Site builder

This Python script takes the data from a Baserow instance as well as the config file and generates a web folder which is then deployed to the given location.

### Installation

```bash
python3 -m venv .venv
source .venv/bin/activate
pip3 install .
```

### Baserow structure

- Table `Bewerbung` (entry)
    - Text Field `Name` (artist name)
    - Single Select Field `Genre`
    - Text Field `Titel`
    - Duration Field `Dauer`
    - Text Field `Beschreibung`
- Table `Spielstätte` 
    - Text Field `Name`
- Table `Timetable` 
    - Datetime Field `Wann`
    - Link Field `Beitrag` (link to registration)
    - Link Field `Spielstätte` (link to location)
    - Bool Field `Permanent?` (is installation or exhibition over the whole festival)

You can alter the names of the fields using the config file.


### Config file

The configuration happens in the files `config.toml` and `secrets.toml`. You find template for `secrets.toml` templates of them in the root of the repository. Create a new project folder with a copy of them and adapt them as needed.

```bash
mkdir my-event
cp config.tpl.toml my-event/config.toml
cp secrets.tpl.toml my-event/secrets.toml
```

You can alter some of the images (logos etc.) displayed. Just set the paths in the config file (if empty default image will be used). The paths are relative to the config file. Example:

```toml
event_name = "My Event – A festival"
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


## Live-reload server for development

There is the possibility to run local webserver with live-reload to speed up the development. Usage:

```bash
pip3 install .[dev]
feti serve /var/www/data
``` 