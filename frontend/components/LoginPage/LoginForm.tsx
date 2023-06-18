import { 
    Box, 
    Stack 
} from '@chakra-ui/react';
import React, { useState } from "react";
import { UsernameInput } from './UsernameInput';
import { PasswordInput } from './PasswordInput';
import { FormSubmit } from './FormSubmit';

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
                <UsernameInput value={form.username} onChange={handleInputChange}/>
                <PasswordInput value={form.password} onChange={handleInputChange}/>
                <FormSubmit username={form.username} password={form.password}/>
            </Stack>
        </Box>
    );
}
