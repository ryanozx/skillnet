import { 
    Button, 
    useToast } 
from '@chakra-ui/react';
import React, { MouseEventHandler } from "react";
import axios from "axios";
import { useRouter } from 'next/router'; 

interface LoginButtonProps {
    username: string;
    password: string;
}

export const LoginButton: React.FC<LoginButtonProps> = ({ username, password }) => {
    const toast = useToast();
    const router = useRouter();

    const onSubmit: MouseEventHandler = () => {
        var form_data = new FormData();
        form_data.append('username', username);
        form_data.append('password', password);
        const base_url = process.env.BACKEND_BASE_URL;
        axios.post(base_url + '/login', form_data, {withCredentials: true})
            .then((res) => {
                toast({
                    title: "Form submission successful.",
                    description: "We've received your form data.",
                    status: "success",
                    duration: 5000,
                    isClosable: true,
                });
                router.push("/feed");
            })
            .catch((error) => {
                toast({
                    title: "An error occurred.",
                    description: error.response.data.error,
                    status: "error",
                    duration: 5000,
                    isClosable: true,
                });
            });
    };

    return (
        <Button 
            data-testid="login-button"
            onClick={onSubmit} 
            colorScheme="teal" 
            size="lg" 
            fontSize="md">
            Sign in
        </Button>
    );
}
