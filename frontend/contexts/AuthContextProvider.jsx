import React, { createContext, useState, useEffect } from 'react';
import axios from 'axios';

export const AuthContext = createContext();

const AuthContextProvider = ({ children }) => {
  const [isLoggedIn, setIsLoggedIn] = useState(true);

    // useEffect(() => {
    //     const token = sessionStorage.getItem('authToken');
    //     console.log("AuthContext checks for token");
        
    //     if (token) {
    //     axios.post('your-backend-url/api/auth', {}, {
    //         headers: {
    //         Authorization: `Bearer ${token}`
    //         }
    //         }).then(response => {
    //             if (response.data.isAuthenticated) {
    //                 setIsLoggedIn(true);
    //             }
    //         }).catch(err => console.error(err));
    //     }
    // }, []);

  return (
    <AuthContext.Provider value={{ isLoggedIn, setIsLoggedIn }}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthContextProvider;
