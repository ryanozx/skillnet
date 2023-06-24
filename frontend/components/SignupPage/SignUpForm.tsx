import React, { useState, MouseEventHandler } from 'react';
import { Stack, Button, FormControl, FormLabel, Input, InputGroup, InputRightElement, useToast } from '@chakra-ui/react';
import LoginRedirect from './LoginRedirect';
import { ViewIcon, ViewOffIcon } from '@chakra-ui/icons';
import axios from 'axios';
import { useRouter } from 'next/router';


export default function SignUpForm() {
    const [showPassword, setShowPassword] = useState<boolean>(false);
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
        const url = 'http://localhost:8080/signup'
    
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
            <FormControl id="username" isRequired>
                <FormLabel>Username</FormLabel>
                <Input type="text" value={username} onChange={e => setUsername(e.target.value)} />
            </FormControl>
            <FormControl id="email" isRequired>
                <FormLabel>Email address</FormLabel>
                <Input type="email" value={email} onChange={e => setEmail(e.target.value)} />
            </FormControl>
            <FormControl id="password" isRequired>
                <FormLabel>Password</FormLabel>
                <InputGroup>
                    <Input type={showPassword ? 'text' : 'password'} value={password} onChange={e => setPassword(e.target.value)} />
                    <InputRightElement h={'full'}>
                        <Button
                            variant={'ghost'}
                            onClick={() =>
                                setShowPassword((showPassword) => !showPassword)
                            }>
                            {showPassword ? <ViewIcon /> : <ViewOffIcon />}
                        </Button>
                    </InputRightElement>
                </InputGroup>
            </FormControl>
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
            <LoginRedirect />
        </Stack>

    );
}