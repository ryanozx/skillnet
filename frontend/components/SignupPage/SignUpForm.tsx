import React, { useState, MouseEventHandler } from 'react';
import { Stack, useToast } from '@chakra-ui/react';
import LoginRedirect from './LoginRedirect';
import { EmailInput } from './EmailInput';
import { UsernameInput } from './UsernameInput';
import { SignUpButton } from './SignUpButton';
import { PasswordInput } from './PasswordInput';

type UserSignupForm = {
    username: string;
    email: string;
    password: string;
};

export default function SignUpForm() {
    const [form, setForm] = useState<UserSignupForm>({ username: "", email: "", password: "" });
    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setForm({
            ...form,
            [name]: value
        });
    };

    return (
        <Stack spacing={4}>
            <UsernameInput data-testid="username" value={form.username} onChange={handleInputChange} />
            <EmailInput data-testid="email" value={form.email} onChange={handleInputChange} />
            <PasswordInput data-testid="password" value={form.password} onChange={handleInputChange} />
            <SignUpButton data-testid="submit" form={form} />
            <LoginRedirect data-testid="redirect" />
        </Stack>
    );
    
}
