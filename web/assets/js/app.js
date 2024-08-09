document.addEventListener("alpine:init", () => {
    var elements = document.querySelectorAll(".sidenav");
    M.Sidenav.init(elements, {});
    Alpine.data("feti", () => ({
        event_name: null,
        event_description: null,
        genre_color_class: null,
        created_on: null,
        permanent: [],
        per_day: null,
        map: null,
        map_loaded: false,
        selected_date: null,
        selected_date_raw: null,
        page_title() {
            if (this.event_name === null) {
                return "Programm – Festival für Freunde";
            }
            return this.event_name;
        },
        async fetch_schedule() {
            try {
                let rsp = await fetch("/schedule.json");
                let data = await rsp.json();
                this.event_name = data.event_name;
                this.event_description = data.event_description;
                this.genre_color_class = data.genre_color_class;
                this.created_on = data.created_on;
                this.permanent = data.permanent;
                this.per_day = data.per_day;
                this.map = data.map;
            } catch (err) {
                console.error(`error fetching schedule: ${err}`)
            }

            // Determines which day should be displayed.
            const today = new Date().toISOString().split("T")[0];
            let dates = Object.keys(this.per_day);
            if (dates.includes(today)) {
                this.selected_date_raw = today;
            } else {
                this.selected_date_raw = dates[0];
            }
            this.selected_date = Date(this.selected_date_raw);
        },
        date_tabs_data() {
            let rsl = {};
            for (let day in this.per_day) {
                rsl[day] = this.render_date(day);
            }
            return rsl;
        },
        render_date(date_str) {
            const day_names = ["So", "Mo", "Di", "Mi", "Do", "Fr", "Sa"];
            const date = new Date(date_str);
            const day_name = day_names[date.getDay()];
            const day_number = date.getDate();
            const month = date.getMonth() + 1;
            return `${day_name}, ${day_number}.${month}`;
        },
        render_long_date(date_str) {
            const day_names = ["Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag"];
            const date = new Date(date_str);
            const day_name = day_names[date.getDay()];
            const day_number = date.getDate();
            const month = date.getMonth() + 1;
            return `${day_name}, ${day_number}.${month}`;
        },
        render_time(date_str) {
            const date = new Date(date_str);
            const hour = date.getHours();
            const minute = String(date.getMinutes()).padStart(2, "0");
            return `${hour}:${minute}`;
        },
        render_date_time(date_str) {
            return `${this.render_long_date(date_str)} ${this.render_time(date_str)}`;
        },
        change_day(day) {
            this.selected_date = new Date(day);
            this.selected_date_raw = day;
        },
        get_genre_color(genre) {
            if (this.genre_color_class.hasOwnProperty(genre)) {
                let rsl = this.genre_color_class[genre].split(" ");
                return rsl;
            }
            return ["is-link"];
        },
        lead(text, word_count) {
            if (word_count <= 0) return "";
            const words = text.split(" ");
            if (words.length <= word_count) {
                return words.join(" ");
            } else {
                return words.slice(0, word_count).join(" ");
            }
        },
        after_lead(text, word_count) {
            if (word_count <= 0) return text;
            const words = text.split(" ");
            if (words.length <= word_count) {
                return "";
            } else {
                return words.slice(word_count).join(" ");
            }
        },
        init_map() {
            let image_bounds = [[0, 0], [this.map.upper_bound_y, this.map.upper_bound_x]];
            let map = L.map("map-container", {
                center: [this.map.center_y, this.map.center_x],
                zoom: this.map.zoom,
                crs: L.CRS.Simple,
                maxBounds: image_bounds,
            });
            L.imageOverlay(
                "assets/img/map.png",
                image_bounds,
            ).addTo(map);
            this.map_loaded = true;
        },
        load_map() {
            if (window.L) {
                this.init_map();
                return
            }
            let script = document.createElement("script");
            script.src = "assets/js/leaflet.js";
            script.async = true;
            script.onload = () => {
                this.init_map();
            }
            document.body.appendChild(script);
        },
    }));
})