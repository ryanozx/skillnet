import React from 'react';
import { Heading, Stack, Text } from '@chakra-ui/react';

export default function FormHeading() {
    return (
        <Stack align={'center'}>
            <Heading data-testid="form-heading" fontSize={'4xl'}>Sign in to your account</Heading>
            <Text data-testid="form-subheading" fontSize={'lg'} color={'gray.600'}>
                to enjoy all of our cool features ✌️
            </Text>
        </Stack>
    );
}
