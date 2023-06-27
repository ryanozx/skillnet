import React from 'react';
import {Stack, Text, Link} from '@chakra-ui/react';

export default function SignUpSection() {
    return (
        <Stack pt={6}>
            <Text align={'center'}>
                Don't have an account? <Link href='/signup' color={'blue.400'}>Sign up</Link>
            </Text>
        </Stack>
    );
}