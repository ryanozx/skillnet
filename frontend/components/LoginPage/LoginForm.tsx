import { 
    Box, 
    Stack 
} from '@chakra-ui/react';
import React, { useState } from "react";
import { UsernameInput } from './UsernameInput';
import { PasswordInput } from './PasswordInput';
import { LoginButton } from './LoginButton';
import SignUpRedirect from './SignUpRedirect';

type UserLoginForm = {
    username: string;
    password: string;
};

export default function LoginForm() {
    const [form, setForm] = useState<UserLoginForm>({ username: "", password: "" });

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setForm((prevForm) => ({ ...prevForm, [name]: value }));
    };

    return (
        <Box
            rounded={'lg'}
            bg='white'
            boxShadow={'lg'}
            p={8}
            w={{ base: '90vw', md: '60vw', lg: '30vw' }}
        >
            <Stack spacing={4}>
                <UsernameInput data-testid="username-input-component" value={form.username} onChange={handleInputChange}/>
                <PasswordInput data-testid="password-input-component" value={form.password} onChange={handleInputChange}/>
                <LoginButton data-testid="login-button" username={form.username} password={form.password}/>
                <SignUpRedirect />
            </Stack>
        </Box>
    );
}
