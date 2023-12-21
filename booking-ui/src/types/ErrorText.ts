import RuntimeConfig from "@/components/RuntimeConfig";
import { TFunction } from "i18next";

//var ResponseCodeBookingSlotConflict: number             = 1001;
var ResponseCodeBookingLocationMaxConcurrent: number = 1002;
var ResponseCodeBookingTooManyUpcomingBookings: number = 1003;
var ResponseCodeBookingTooManyDaysInAdvance: number = 1004;
var ResponseCodeBookingInvalidBookingDuration: number = 1005;
var ResponseCodeBookingMaxConcurrentForUser: number = 1006;

export default class ErrorText {
    static getTextForAppCode(code: number, t: TFunction): string {
        if (code === ResponseCodeBookingLocationMaxConcurrent) {
            return t("errorTooManyConcurrent");
        } else if (code === ResponseCodeBookingInvalidBookingDuration) {
            return t("errorBookingDuration", { "num": RuntimeConfig.INFOS.maxBookingDurationHours });
        } else if (code === ResponseCodeBookingTooManyDaysInAdvance) {
            return t("errorDaysAdvance", { "num": RuntimeConfig.INFOS.maxDaysInAdvance });
        } else if (code === ResponseCodeBookingTooManyUpcomingBookings) {
            return t("errorBookingLimit", { "num": RuntimeConfig.INFOS.maxBookingsPerUser });
        } else if (code === ResponseCodeBookingMaxConcurrentForUser) {
            return t("errorConcurrentBookingLimit", { "num": RuntimeConfig.INFOS.maxConcurrentBookingsPerUser });
        } else {
            return t("errorUnknown");
        }
    }
}
