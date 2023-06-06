import React from 'react';
import { Heading, Stack, Text } from '@chakra-ui/react';

export default function FormHeading() {
    return (
        <Stack align={'center'}>
            <Heading fontSize={'4xl'} textAlign={'center'}>
                Sign up
            </Heading>
            <Text fontSize={'lg'} color={'gray.600'}>
                to enjoy all of our cool features ✌️
            </Text>
        </Stack>
    );
}