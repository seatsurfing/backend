import React from 'react';

export interface AuthContextData {
    token: string;
    username: string;
    isLoading: boolean;
    maxBookingsPerUser: number;
    maxDaysInAdvance: number;
    maxBookingDurationHours: number;
    dailyBasisBooking: boolean;
    showNames: boolean;
    setDetails: (token: string, username: string) => void;
};

export const AuthContext = React.createContext<AuthContextData>({
    token: "", 
    username: "", 
    isLoading: true, 
    maxBookingsPerUser: 0,
    maxDaysInAdvance: 0,
    maxBookingDurationHours: 0,
    dailyBasisBooking: false,
    showNames: false,
    setDetails: (token: string, username: string) => {},
});
