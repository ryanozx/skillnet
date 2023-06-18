import React, { createContext, useState, useEffect, useContext } from 'react';
import axios from 'axios';

interface UserContextProps {
    needUpdate: boolean;
    setNeedUpdate: React.Dispatch<React.SetStateAction<boolean>>;
}

interface UserProviderProps {
    children: React.ReactNode;
}

const UserContext = createContext<UserContextProps | undefined>(undefined);

export const UserProvider: React.FC<UserProviderProps> = ({ children }) => {
    const [needUpdate, setNeedUpdate] = useState<boolean>(true);
    return (
        <UserContext.Provider value={{ needUpdate, setNeedUpdate }}>
            {children}
        </UserContext.Provider>
    );
};


export const useUser = () => {
    const context = useContext(UserContext);
    if (context === undefined) {
        throw new Error('useUser must be used within a UserProvider');
    }
    return context;
};
