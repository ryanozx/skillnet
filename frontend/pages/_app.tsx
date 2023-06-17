import { ChakraProvider } from '@chakra-ui/react';
import { Provider } from 'react-redux';
import React from 'react';

interface MyAppProps {
  Component: React.ComponentType<any>;
  pageProps: any;
}

const MyApp: React.FC<MyAppProps> = ({ Component, pageProps }) => {
  return (
    <React.StrictMode>
      <ChakraProvider>

          <Component {...pageProps} />
      </ChakraProvider>
    </React.StrictMode>
  );
};

export default MyApp;
