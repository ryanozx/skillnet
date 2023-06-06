import React from 'react';
import { Button, Link, Stack } from '@chakra-ui/react';

export const CallToActionButtons = () => (
    <Stack
        direction={'column'}
        spacing={3}
        align={'center'}
        alignSelf={'center'}
        position={'relative'}>
        <Link href='/login'>
            <Button
                colorScheme={'green'}
                bg={'green.400'}
                rounded={'full'}
                px={6}
                _hover={{
                    bg: 'green.500',
                }}>
                Log In
            </Button>
        </Link>
        <Link href='/signup'>
            <Button variant={'link'} colorScheme={'blue'} size={'sm'}>
                Sign Up
            </Button>
        </Link>
    </Stack>
);