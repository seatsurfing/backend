export default class Formatting {
    static getFormatter(): Intl.DateTimeFormat {
        let formatter = new Intl.DateTimeFormat('de', {
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

    static getFormatterShort(): Intl.DateTimeFormat {
        let formatter = new Intl.DateTimeFormat('de', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: 'numeric',
            minute: 'numeric',
            hour12: false
          });
        return formatter;
    }

    static getFormatterDate(): Intl.DateTimeFormat {
        let formatter = new Intl.DateTimeFormat('de', {
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
            return "Heute";
        }
        if (start == today+1) {
            return "Morgen";
        }
        if (start == today+2) {
            return "Übermorgen";
        }
        if (start > today && start <= today+7) {
            return "In " + (start-today) + " Tagen";
        }
        return Formatting.getFormatterDate().format(enter);
    }
}