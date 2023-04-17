import React, { useState } from 'react';

/*
export interface AuthContextData {
    username: string;
    isLoading: boolean;
    maxBookingsPerUser: number;
    maxDaysInAdvance: number;
    maxBookingDurationHours: number;
    dailyBasisBooking: boolean;
    showNames: boolean;
    defaultTimezone: string;
    //setDetails: (username: string) => void;
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
    //setDetails: (username: string) => { },
});

interface Props {
    children: JSX.Element
}

export const AuthContextProvider = (props: Props) => {
    const [username, setUsername] = useState("");
    const [isLoading, setIsLoading] = useState(true);
    const [maxBookingsPerUser, setMaxBookingsPerUser] = useState(0);
    const [maxDaysInAdvance, setMaxDaysInAdvance] = useState(0);
    const [maxBookingDurationHours, setMaxBookingDurationHours] = useState(0);
    const [dailyBasisBooking, setDailyBasisBooking] = useState(false);
    const [showNames, setShowNames] = useState(false);
    const [defaultTimezone, setDefaultTimezone] = useState("");
    //const [setDetails, setSetDetails] = useState((string) => void);

    return (
        <AuthContext.Provider
            value={{
                username: username,
                isLoading: isLoading,
                maxBookingsPerUser: maxBookingsPerUser,
                maxDaysInAdvance: maxDaysInAdvance,
                maxBookingDurationHours: maxBookingDurationHours,
                dailyBasisBooking: dailyBasisBooking,
                showNames: showNames,
                defaultTimezone: defaultTimezone,
                //setDetails: setDetails,
            }}
        >
            {props.children}
        </AuthContext.Provider>
    );
};
*/