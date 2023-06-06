import React from 'react';
import { Flex, Stack, useColorModeValue } from '@chakra-ui/react';
import FormHeading from './FormHeading';
import LoginForm from './LoginForm';

export default function LoginPageContainer() {
    return (
        <Flex
            minH={'100vh'}
            align={'center'}
            justify={'center'}
            bg={useColorModeValue('gray.50', 'gray.800')}
        >
            <Stack spacing={8} mx={'auto'} py={12} px={6}>
                <FormHeading />
                <LoginForm />
            </Stack>
        </Flex>
    );
}