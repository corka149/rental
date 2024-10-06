const locales = {
  en: {
    days: [
      "Sunday",
      "Monday",
      "Tuesday",
      "Wednesday",
      "Thursday",
      "Friday",
      "Saturday",
    ],
    daysShort: ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"],
    daysMin: ["Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"],
    months: [
      "January",
      "February",
      "March",
      "April",
      "May",
      "June",
      "July",
      "August",
      "September",
      "October",
      "November",
      "December",
    ],
    monthsShort: [
      "Jan",
      "Feb",
      "Mar",
      "Apr",
      "May",
      "Jun",
      "Jul",
      "Aug",
      "Sep",
      "Oct",
      "Nov",
      "Dec",
    ],
    today: "Today",
    clear: "Clear",
    dateFormat: "MM/dd/yyyy",
    timeFormat: "hh:mm aa",
    firstDay: 0,
  },
  de: {
    days: [
      "Sonntag",
      "Montag",
      "Dienstag",
      "Mittwoch",
      "Donnerstag",
      "Freitag",
      "Samstag",
    ],
    daysShort: ["Son", "Mon", "Die", "Mit", "Don", "Fre", "Sam"],
    daysMin: ["So", "Mo", "Di", "Mi", "Do", "Fr", "Sa"],
    months: [
      "Januar",
      "Februar",
      "März",
      "April",
      "Mai",
      "Juni",
      "Juli",
      "August",
      "September",
      "Oktober",
      "November",
      "Dezember",
    ],
    monthsShort: [
      "Jan",
      "Feb",
      "Mär",
      "Apr",
      "Mai",
      "Jun",
      "Jul",
      "Aug",
      "Sep",
      "Okt",
      "Nov",
      "Dez",
    ],
    today: "Heute",
    clear: "Löschen",
    dateFormat: "dd.MM.yyyy",
    timeFormat: "HH:mm",
    firstDay: 1,
  },
};

const controller = {
  currentObject: null,
  entries: [],
  getEntries: async function (month, year) {

    month = month + 1;
    
    const searchParams = new URLSearchParams({
      month: month,
      year: year,
      object: this.currentObject ?? "1",
    });

    const response = await fetch("/calendar/search?" + searchParams.toString());

    const entries = await response.json();

    this.entries = (entries ?? []).map(({occurs_on}) => new Date(occurs_on));
  },
  onChangeViewDate: async function ({ month, year, decade }) {
    await this.getEntries(month, year);
  },
  onRenderCell: function ({ date, cellType }) {
    if (cellType === "day") {
      let notAvailable = false;

      for (const entry of this.entries) {
        if (notAvailable) {
          break;
        }

        notAvailable = 
          (date.getFullYear() === entry.getFullYear() &&
          date.getMonth() === entry.getMonth() &&
          date.getDate() === entry.getDate());
      }

      return {
        html: notAvailable ? "<span class='text-red-600'>X</span>" : undefined,
        disabled: notAvailable,
      };
    }
  },
  newCalendar: async function () {
    const lang = document.getElementsByTagName("html")[0].lang;
    const locale = locales[lang] ?? locales.en;

    const now = new Date();
    const month = now.getMonth();
    const year = now.getFullYear();

    await this.getEntries(month, year);

    this.calendar = new AirDatepicker("#calendar", {
      inline: true,
      locale: locale,
      onChangeViewDate: this.onChangeViewDate.bind(this),
      onRenderCell: this.onRenderCell.bind(this),
    });
  },
  registerForObjectSelect: function() {
    var objectSelect = document.getElementById("object");

    this.currentObject = objectSelect.value;

    objectSelect.addEventListener("change", async (event) => {
      if (this.currentObject !== event.target.value) {
        this.currentObject = event.target.value;

        if (this.calendar) {
          this.calendar.destroy();
        }

        this.newCalendar();
      }
    });
  },
};

document.addEventListener("DOMContentLoaded", async () => {
  controller.registerForObjectSelect();
  controller.newCalendar();
});
