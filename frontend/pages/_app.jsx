import { ChakraProvider } from '@chakra-ui/react'
import React from 'react'
import App from 'next/app'

function MyApp({ Component, pageProps }) {
  return (
    <ChakraProvider>
      <Component {...pageProps} />
    </ChakraProvider>
  )
}

export default MyApp;