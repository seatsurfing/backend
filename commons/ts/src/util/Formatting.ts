export default class Formatting {
    static Language: string = "de";
    static I18n = {
        "de": {
            "today": "Heute",
            "tomorrow": "Morgen",
            "inXdays": "In {{x}} Tagen",
        },
        "en": {
            "today": "Today",
            "tomorrow": "Tomorrow",
            "inXdays": "In {{x}} days",
        },
    };

    static t(s: string, vars?: { [key: string]: any }) {
        let res: string = Formatting.I18n[Formatting.Language][s];
        if (!res) {
            return s;
        }
        if (vars) {
            for (const k in vars) {
                res = res.replaceAll("{{"+k+"}}", vars[k]);
            }
        }
        return res;
    }

    static getFormatter(local?: boolean): Intl.DateTimeFormat {
        let formatter = new Intl.DateTimeFormat(Formatting.Language, {
            timeZone: local ? undefined : 'UTC',
            weekday: 'long',
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: 'numeric',
            minute: 'numeric',
            hour12: false
          });
        return formatter;
    }

    static getFormatterNoTime(local?: boolean): Intl.DateTimeFormat {
        let formatter = new Intl.DateTimeFormat(Formatting.Language, {
            timeZone: local ? undefined : 'UTC',
            weekday: 'long',
            year: 'numeric',
            month: '2-digit',
            day: '2-digit'
          });
        return formatter;
    }

    static getFormatterShort(local?: boolean): Intl.DateTimeFormat {
        let formatter = new Intl.DateTimeFormat(Formatting.Language, {
            timeZone: local ? undefined : 'UTC',
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: 'numeric',
            minute: 'numeric',
            hour12: false
          });
        return formatter;
    }

    static getFormatterDate(local?: boolean): Intl.DateTimeFormat {
        let formatter = new Intl.DateTimeFormat(Formatting.Language, {
            timeZone: local ? undefined : 'UTC',
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
          });
        return formatter;
    }

    static getDayValue(date: Date): number {
        let s = date.getFullYear().toString().padStart(4, "0") + (date.getMonth()+1).toString().padStart(2, "0") + date.getDate().toString().padStart(2, "0");
        return parseInt(s);
    }

    static getISO8601(date: Date): string {
        let s = date.getFullYear().toString().padStart(4, "0") + "-" + (date.getMonth()+1).toString().padStart(2, "0") + "-" + date.getDate().toString().padStart(2, "0");
        return s;
    }

    static getDateOffsetText(enter: Date, leave: Date): string {
        let today = Formatting.getDayValue(new Date());
        let start = Formatting.getDayValue(enter);
        let end = Formatting.getDayValue(leave);
        if (start <= today && today <= end) {
            return Formatting.t("today");
        }
        if (start == today+1) {
            return Formatting.t("tomorrow");
        }
        if (start > today && start <= today+7) {
            return Formatting.t("inXdays", {"x": (start-today)});
        }
        return Formatting.getFormatterDate().format(enter);
    }

    static convertToFakeUTCDate(d: Date): Date {
        return new Date(Date.UTC(d.getFullYear(), d.getMonth(), d.getDate(), d.getHours(), d.getMinutes(), d.getSeconds(), 0));
    }

    static stripTimezoneDetails(s: string): string {
        if ((s.length > 6) && ((s[s.length-6] === "+") || (s[s.length-6] === "-"))) {
            return s.substring(0, s.length-6) + ".000Z";
        }
        return s;
    }
}