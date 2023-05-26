import { ChakraProvider } from '@chakra-ui/react';
import { Provider } from 'react-redux';
import store from '../store';
import React from 'react';

export default function MyApp({ Component, pageProps }) {
  return (
    <React.StrictMode>
        <ChakraProvider>
            <Provider store={store}>
                <Component {...pageProps} />
            </Provider>
        </ChakraProvider>

    </React.StrictMode>
    
  );
}
