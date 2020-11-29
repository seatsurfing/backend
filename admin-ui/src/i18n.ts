import i18n from 'i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import { initReactI18next } from 'react-i18next';

const resources = {
    de: {
        translation: {
            "weweaveUG": "weweave UG (haftungsbeschränkt)",
            "privacy": "Datenschutz",
            "imprint": "Impressum",
            "mangageOrgHeadline": "Organisation verwalten.",
            "emailAddress": "E-Mail Adresse",
            "errorInvalidEmail": "Ungültige E-Mail-Adresse.",
            "errorNoAuthProviders": "Für diesen Nutzer stehen keine Anmelde-Möglichkeiten zur Verfügung.",
            "errorInvalidPassword": "Ungültiges Kennwort.",
            "password": "Kennwort",
            "signin": "Anmelden",
            "back": "Zurück",
            "signinAsAt": "Als {{user}} an {{org}} anmelden:",
        }
    },
    en: {
        translation: {
            "weweaveUG": "weweave UG (limited liability)",
            "privacy": "Privacy",
            "imprint": "Imprint",
            "mangageOrgHeadline": "Manage organization.",
            "emailAddress": "Email address",
            "errorInvalidEmail": "Invalid email address.",
            "errorNoAuthProviders": "No authentication providers available for this user.",
            "errorInvalidPassword": "Invalid password.",
            "password": "Password",
            "signin": "Sign in",
            "back": "Back",
            "signinAsAt": "Sign is as {{user}} at {{org}}:",
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

export default i18n;
