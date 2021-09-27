import { TFunction } from "i18next";

//var ResponseCodeBookingSlotConflict: number             = 1001;
var ResponseCodeBookingLocationMaxConcurrent: number = 1002;
var ResponseCodeBookingTooManyUpcomingBookings: number = 1003;
var ResponseCodeBookingTooManyDaysInAdvance: number = 1004;
var ResponseCodeBookingInvalidBookingDuration: number = 1005;

export default class ErrorText {
    static getTextForAppCode(code: number, t: TFunction, context: any): string {
        if (code === ResponseCodeBookingLocationMaxConcurrent) {
            return t("errorTooManyConcurrent");
        } else if (code === ResponseCodeBookingInvalidBookingDuration) {
            return t("errorBookingDuration", { "num": context.maxBookingDurationHours });
        } else if (code === ResponseCodeBookingTooManyDaysInAdvance) {
            return t("errorDaysAdvance", { "num": context.maxDaysInAdvance });
        } else if (code === ResponseCodeBookingTooManyUpcomingBookings) {
            return t("errorBookingLimit", { "num": context.maxBookingsPerUser });
        } else {
            return t("errorUnknown");
        }
    }
}
