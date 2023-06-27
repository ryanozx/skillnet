import React, { MouseEventHandler } from 'react';
import { Stack, Button, FormControl, FormLabel, Input, InputGroup, InputRightElement, useToast } from '@chakra-ui/react';
import LoginRedirect from './LoginRedirect';
import { ViewIcon, ViewOffIcon } from '@chakra-ui/icons';
import axios from 'axios';
import { useRouter } from 'next/router';

interface SubmitButtonProps {
    onClick: MouseEventHandler<HTMLButtonElement>;
}

export const SubmitButton = ({ onClick }: SubmitButtonProps) => (
    <Stack spacing={10} pt={2}>
        <Button
            type="submit"
            loadingText="Submitting"
            size="lg"
            bg={'blue.400'}
            color={'white'}
            _hover={{
                bg: 'blue.500',
            }}
            onClick={onClick}>
            Sign up
        </Button>
    </Stack>
);
