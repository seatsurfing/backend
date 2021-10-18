import React from 'react';

export interface AuthContextData {
    username: string;
    isLoading: boolean;
    maxBookingsPerUser: number;
    maxDaysInAdvance: number;
    maxBookingDurationHours: number;
    dailyBasisBooking: boolean;
    showNames: boolean;
    defaultTimezone: string;
    setDetails: (username: string) => void;
};

export const AuthContext = React.createContext<AuthContextData>({
    username: "", 
    isLoading: true, 
    maxBookingsPerUser: 0,
    maxDaysInAdvance: 0,
    maxBookingDurationHours: 0,
    dailyBasisBooking: false,
    showNames: false,
    defaultTimezone: "",
    setDetails: (username: string) => {},
});
