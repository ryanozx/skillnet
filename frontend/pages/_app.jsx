import { ChakraProvider } from '@chakra-ui/react';
import { Provider } from 'react-redux';
import store from '../store';

export default function MyApp({ Component, pageProps }) {
  return (
    <ChakraProvider>
        <Provider store={store}>
            <Component {...pageProps} />
        </Provider>
    </ChakraProvider>
  );
}
