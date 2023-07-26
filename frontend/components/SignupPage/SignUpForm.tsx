import React, { useState, useEffect, useRef } from 'react';
import { Stack } from '@chakra-ui/react';
import { EmailInput } from './EmailInput';
import { UsernameInput } from './UsernameInput';
import { SignUpButton } from './SignUpButton';
import { PasswordInput } from './PasswordInput';
import { RetypePasswordInput } from './RetypePasswordInput';
import { useToast } from '@chakra-ui/react';
import { useRouter } from 'next/router';
import axios from 'axios';
import { escapeHtml } from '../../types';

type UserSignupForm = {
    username: string;
    email: string;
    password: string;
    retypePassword: string;
};

export default function SignUpForm() {
    const [form, setForm] = useState<UserSignupForm>({ username: "", email: "", password: "", retypePassword: "" });

    const [usernameChanged, setUsernameChanged] = useState<boolean>(false);
    const [usernameError, setUsernameError] = useState<boolean>(false);

    const [emailChanged, setEmailChanged] = useState<boolean>(false);
    const [emailError, setEmailError] = useState<boolean>(false);

    const [passwordChanged, setPasswordChanged] = useState<boolean>(false);
    const [passwordError, setPasswordError] = useState<boolean>(false);

    const [retypePasswordChanged, setRetypePasswordChanged] = useState<boolean>(false);
    const [retypePasswordError, setRetypePasswordError] = useState<boolean>(false);

    const [passwordMismatch, setPasswordMismatch] = useState<boolean>(false);

    const [allFieldsChanged, setAllFieldsChanged] = useState<boolean>(false);
    const [formNoError, setFormNoError] = useState<boolean>(false);

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setForm({
            ...form,
            [name]: value
        });
    };

    useEffect(() => setAllFieldsChanged(usernameChanged && emailChanged && passwordChanged && retypePasswordChanged), 
    [usernameChanged, emailChanged, passwordChanged, retypePasswordChanged]);

    useEffect(() => setFormNoError(
        !(usernameError || emailError || passwordError || retypePasswordError)
    ), [usernameError, emailError, passwordError, retypePasswordError])

    useEffect(() => setPasswordMismatch(form.password != form.retypePassword), [form]);

    const toast = useToast();
    const router = useRouter();

    const handleSubmit = async (e : React.FormEvent<HTMLFormElement>) => {    
        e.preventDefault();    
        var form_data = new FormData();
        form_data.append('email', form.email);
        form_data.append('username', escapeHtml(form.username));
        form_data.append('password', form.password);
        const base_url = process.env.BACKEND_BASE_URL;
        const url = base_url + '/signup';
    
        axios.post(url, form_data, {withCredentials: true})
            .then(() => {
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
                    title: "Signup unsuccessful",
                    description: error.response.data.error,
                    status: "error",
                    duration: 5000,
                    isClosable: true,
                });
            });
    }

    return (
        <Stack spacing={4}>
            <form onSubmit={e => allFieldsChanged && formNoError && handleSubmit(e)}>
                <UsernameInput data-testid="username" value={form.username} onChange={handleInputChange} 
                usernameChanged={usernameChanged} setUsernameChanged={setUsernameChanged} setUsernameError={setUsernameError}/>
                <EmailInput data-testid="email" value={form.email} onChange={handleInputChange}
                emailChanged={emailChanged} setEmailChanged={setEmailChanged} setEmailError={setEmailError}/>
                <PasswordInput data-testid="password" value={form.password} onChange={handleInputChange}
                passwordChanged={passwordChanged} setPasswordChanged={setPasswordChanged} setPasswordError={setPasswordError}/>
                <RetypePasswordInput data-testid="retypePassword" value={form.retypePassword} onChange={handleInputChange} 
                passwordMismatch={passwordMismatch} setRetypePasswordError={setRetypePasswordError}
                retypePasswordChanged={retypePasswordChanged} setRetypePasswordChanged={setRetypePasswordChanged}/>
                <SignUpButton data-testid="submit" formValid={allFieldsChanged && formNoError}/>
            </form>
        </Stack>
    );
    
}
