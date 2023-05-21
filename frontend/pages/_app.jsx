import { ChakraProvider } from '@chakra-ui/react';
import AuthContextProvider from '../contexts/AuthContextProvider.jsx';

export default function MyApp({ Component, pageProps }) {
  return (
    <ChakraProvider>
        <AuthContextProvider>
            <Component {...pageProps} />
        </AuthContextProvider>
    </ChakraProvider>
  )
}
