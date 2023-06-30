import React from 'react';
import {Stack, Text, Link} from '@chakra-ui/react';

export default function SignUpRedirect() {
    return (
        <Stack pt={6}>
            <Text data-testid="redirect-text" align={'center'}>
                Don't have an account? <Link data-testid="signup-link" href='/signup' color={'blue.400'}>Sign up</Link>
            </Text>
        </Stack>
    );
}