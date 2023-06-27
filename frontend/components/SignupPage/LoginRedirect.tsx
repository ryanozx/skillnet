import React from 'react';
import { Stack, Text, Link } from '@chakra-ui/react';

export default function LoginRedirect() {
    return (
        <Stack pt={6}>
            <Text data-testid="redirect-text" align={'center'}>
                Already a user? <Link data-testid="login-link" href='/login' color={'blue.400'}>Login</Link>
            </Text>
        </Stack>
    );
}
