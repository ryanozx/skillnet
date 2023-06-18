import { ChakraProvider } from '@chakra-ui/react';
import React, { createContext, useState, useEffect, useContext } from 'react';
import axios from 'axios';
import { UserProvider } from '../userContext';

interface MyAppProps {
  Component: React.ComponentType<any>;
  pageProps: any;
}

const MyApp: React.FC<MyAppProps> = ({ Component, pageProps }) => {
  return (
    <React.StrictMode>
      <ChakraProvider>
        <UserProvider>
          <Component {...pageProps} />
        </UserProvider>
      </ChakraProvider>
    </React.StrictMode>
  );
};

export default MyApp;
