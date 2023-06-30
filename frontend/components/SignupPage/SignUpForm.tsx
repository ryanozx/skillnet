import React, { useState, MouseEventHandler } from 'react';
import { Stack, useToast } from '@chakra-ui/react';
import LoginRedirect from './LoginRedirect';
import axios from 'axios';
import { useRouter } from 'next/router';
import { EmailInput } from './EmailInput';
import { UsernameInput } from './UsernameInput';
import { SubmitButton } from './SignUpButton';
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
            <UsernameInput data-testid="username-input" value={form.username} onChange={handleInputChange} />
            <EmailInput data-testid="email-input" value={form.email} onChange={handleInputChange} />
            <PasswordInput data-testid="password-input" value={form.password} onChange={handleInputChange} />
            <SubmitButton data-testid="submit-button" form={form} />
            <LoginRedirect />
        </Stack>
    );
    
}
