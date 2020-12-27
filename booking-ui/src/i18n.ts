import { Formatting } from 'flexspace-commons';
import i18n from 'i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import { initReactI18next } from 'react-i18next';

const resources = {
    de: {
        translation: {
            "weweaveUG": "weweave UG (haftungsbeschränkt)",
            "privacy": "Datenschutz",
            "imprint": "Impressum",
            "findYourPlace": "Finde Deinen Platz.",
            "emailPlaceholder": "max@mustermann.de",
            "getStarted": "Loslegen",
            "back": "Zurück",
            "errorInvalidEmail": "Ungültige E-Mail-Adresse.",
            "errorLogin": "Fehler bei der Anmeldung. Möglicherweise ist diese E-Mail-Adresse nicht mit einer Organisation verknüpft.",
            "errorNoAuthProviders": "Für diesen Nutzer stehen keine Anmelde-Möglichkeiten zur Verfügung.",
            "errorInvalidPassword": "Ungültiges Kennwort.",
            "password": "Kennwort",
            "signin": "Anmelden",
            "signinAsAt": "Als {{user}} an {{org}} anmelden:",
            "enter": "Beginn", 
            "leave": "Ende",
            "area": "Bereich",
            "errorBookingLimit": "Das Limit von {{num}} Buchungen wurde erreicht.",
            "errorPickArea": "Bitte einen Bereich auswählen.",
            "errorEnterFuture": "Der Beginn muss in der Zukunft liegen.",
            "errorLeaveAfterEnter": "Das Ende muss nach dem Beginn liegen.",
            "errorDaysAdvance": "Die Buchung darf maximal {{num}} Tage in der Zukunft liegen.",
            "errorBookingDuration": "Die maximale Buchungsdauer beträgt {{num}} Stunden.",
            "searchSpace": "Plätze suchen",
            "pleaseSelect": "bitte wählen",
            "signout": "Abmelden",
            "bookSeat": "Platz buchen",
            "myBookings": "Meine Buchungen",
            "loadingHint": "Daten werden geladen...",
            "space": "Platz",
            "confirmBooking": "Buchung bestätigen",
            "cancel": "Abbrechen",
            "bookingConfirmed": "Deine Buchung wurde bestätigt!",
            "cancelBooking": "Buchung stornieren",
            "confirmCancelBooking": "Buchung für {{enter}} stornieren?",
            "noBookings": "Keine Buchungen gefunden.",
        }
    },
    en: {
        translation: {
            "weweaveUG": "weweave UG (limited liability)",
            "privacy": "Privacy",
            "imprint": "Imprint",
            "findYourPlace": "Find your space.",
            "emailPlaceholder": "you@company.com",
            "getStarted": "Get started",
            "back": "Back",
            "errorInvalidEmail": "Invalid email address.",
            "errorLogin": "An error occurred while signing you in. Your email address might not be associated with an organization.",
            "errorNoAuthProviders": "No authentication providers for your user.",
            "errorInvalidPassword": "Invalid password.",
            "password": "Password",
            "signin": "Sign in",
            "signinAsAt": "Sign is as {{user}} at {{org}}:",
            "enter": "Enter", 
            "leave": "Leave",
            "area": "Area",
            "errorBookingLimit": "You've reached the limit of {{num}} bookings.",
            "errorPickArea": "Please pick an area.",
            "errorEnterFuture": "Enter date must be in the future.",
            "errorLeaveAfterEnter": "Leave date must be after enter date.",
            "errorDaysAdvance": "Your booking must not be more than {{num}} days in advance.",
            "errorBookingDuration": "The maximum booking duration is {{num}} hours.",
            "searchSpace": "Find a space",
            "pleaseSelect": "please choose",
            "signout": "Sign off",
            "bookSeat": "Book a space",
            "myBookings": "My bookings",
            "loadingHint": "Loading data...",
            "space": "Space",
            "confirmBooking": "Confirm booking",
            "cancel": "Cancel",
            "bookingConfirmed": "Your booking has been confirmed!",
            "cancelBooking": "Cancel booking",
            "confirmCancelBooking": "Cancel your upcoming booking for {{enter}}?",
            "noBookings": "No bookings.",
        }
    }
};

i18n
.use(LanguageDetector)
.use(initReactI18next)
.init({
    resources,
    fallbackLng: "en",
    keySeparator: false
});
Formatting.Language = i18n.language.split("-")[0];

export default i18n;
