import { Box, Stack, FormControl, FormLabel, Input, Button, useColorModeValue, useToast } from '@chakra-ui/react';
import ActionSection from './ActionSection';
import { useState } from "react";
import axios from "axios";


export default function LoginForm() {

    const [form, setForm] = useState({ username: "", password: "" });
    const toast = useToast();

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        setForm((prevForm) => ({ ...prevForm, [name]: value }));
    };

    const onSubmit = () => {
        const {username, password} = form
        var form_data = new FormData();
        form_data.append('username', username);
        form_data.append('password', password);
        axios.post('http://localhost:8080/login', form_data)
            .then((res) => {
                console.log(res);
                toast({
                title: "Form submission successful.",
                description: "We've received your form data.",
                status: "success",
                duration: 5000,
                isClosable: true,
                });
            })
            .catch((error) => {
                console.log(error);
                toast({
                title: "An error occurred.",
                description: error,
                status: "error",
                duration: 5000,
                isClosable: true,
                });
            });
    };

    return (
        <Box
            rounded={'lg'}
            bg={useColorModeValue('white', 'gray.700')}
            boxShadow={'lg'}
            p={8}
            w={{ base: '90vw', md: '60vw', lg: '30vw' }}
        >
            <Stack spacing={4}>
                <FormControl id="username">
                    <FormLabel>username</FormLabel>
                    <Input
                        type="username"
                        name="username"
                        value={form.username}
                        onChange={handleInputChange}
                    />
                </FormControl>
                <FormControl id="password">
                    <FormLabel>Password</FormLabel>
                    <Input
                        type="password"
                        name="password"
                        value={form.password}
                        onChange={handleInputChange}
                    />
                </FormControl>
                <ActionSection onSubmit={onSubmit}/>
            </Stack>
        </Box>
    );
}