document.addEventListener("alpine:init", () => {
    Alpine.data("feti", () => ({
        event_name: null,
        permanent: [],
        per_day: null,
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
                this.permanent = data.permanent;
                // Add expanded property.
                this.per_day = data.per_day;
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
        change_day(day) {
            this.selected_date = new Date(day);
            this.selected_date_raw = day;
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
        }
    }));
})