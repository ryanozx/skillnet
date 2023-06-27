import React, { useState, MouseEventHandler } from 'react';
import { Stack, useToast } from '@chakra-ui/react';
import LoginRedirect from './LoginRedirect';
import axios from 'axios';
import { useRouter } from 'next/router';
import { EmailInput } from './EmailInput';
import { UsernameInput } from './UsernameInput';
import { SubmitButton } from './SubmitButton';
import { PasswordInput } from './PasswordInput';

export default function SignUpForm() {
    const [username, setUsername] = useState<string>('');
    const [email, setEmail] = useState<string>('');
    const [password, setPassword] = useState<string>('');
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
                console.log(error);
                toast({
                    title: "An error occurred.",
                    description: error.response.data.message,
                    status: "error",
                    duration: 5000,
                    isClosable: true,
                });
            });
    }

    return (
        <Stack spacing={4}>
            <UsernameInput data-testid="username-input" value={username} onChange={setUsername} />
            <EmailInput data-testid="email-input" value={email} onChange={setEmail} />
            <PasswordInput data-testid="password-input" value={password} onChange={setPassword} />
            <SubmitButton data-testid="submit-button" onClick={handleClick} />
            <LoginRedirect />
        </Stack>
    );
    
}
