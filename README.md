# feti-ng

HTML5 client side festival timetable. Uses AlpineJS to render the `schedule.json` as a responsive timetable. Can be deployed using any simple HTTP server.

The python script is used to generate a `schedule.json` from our Baserow instance.


## Usage

Use your favorite HTTP-server to serve the content of the `web` directory. Done.


## Generate the schedule

This Python script takes the data from a Baserow instance and generates a `schedule.json`.

### Installation

```bash
python3 -m venv .venv
source .venv/bin/activate
pip3 install .
```

### Baserow structure


- Table `Bewerbung`
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

The configuration happens in the files `config.toml` and `secrets.toml`. You find templates of them in the root of the repository. Create a copy of them and adapt them as needed.

```bash
cp config.tpl.toml config.toml
cp secrets.tpl.toml secrets.toml
```


### Usage 

```bash
feti -c config.toml -s secrets.toml -o schedule.json
```


## Live-reload server for development

There is the possibility to run local webserver with live-reload to speed up the development. Usage:

```bash
pip3 install .[dev]
feti serve
``` 