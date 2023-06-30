import React, { MouseEventHandler } from 'react';
import { Stack, Button, useToast } from '@chakra-ui/react';
import axios from 'axios';
import { useRouter } from 'next/router';

type UserSignupForm = {
    username: string;
    email: string;
    password: string;
};
interface SubmitButtonProps {
    form: UserSignupForm;
}

export const SubmitButton = ({ form }: SubmitButtonProps) => {
    const { username, email, password } = form;
    const toast = useToast();
    const router = useRouter();

    const handleClick: MouseEventHandler = () => {        
        var form_data = new FormData();
        form_data.append('email', email);
        form_data.append('username', username);
        form_data.append('password', password);
        const base_url = process.env.BACKEND_BASE_URL;
        const url = base_url + '/signup';
    
        axios.post(url, form_data, {withCredentials: true})
            .then((res) => {
                router.push('/create-profile');
                toast({
                    title: "Form submission successful.",
                    description: "Account successfully created.",
                    status: "success",
                    duration: 5000,
                    isClosable: true,
                });
            })
            .catch((error) => {
                toast({
                    title: "An error occurred.",
                    description: error.response.data.message,
                    status: "error",
                    duration: 5000,
                    isClosable: true,
                });
            });
    }

    return(
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
                onClick={handleClick}>
                Sign up
            </Button>
        </Stack>
    )
};
