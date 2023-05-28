import { useState } from 'react';
import { Stack, Button, FormControl, FormLabel, Input, InputGroup, InputRightElement, useToast } from '@chakra-ui/react';
import LoginRedirect from './LoginRedirect';
import { ViewIcon, ViewOffIcon } from '@chakra-ui/icons';
import axios from 'axios';

export default function SignUpForm() {
    const [showPassword, setShowPassword] = useState(false);
    const [username, setUsername] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const toast = useToast();

    const handleSubmit = (e) => {
        e.preventDefault();
        
        axios.post('fake-endpoint', {
            username,
            email,
            password
        }).then(() => {
            toast({
                title: "Account created.",
                description: "We've created your account.",
                status: "success",
                duration: 5000,
                isClosable: true,
            });
        }).catch((error) => {
            toast({
                title: "An error occurred.",
                description: error.message,
                status: "error",
                duration: 5000,
                isClosable: true,
            });
        });
    }

    return (
        <form onSubmit={handleSubmit}>
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
                        }}>
                        Sign up
                    </Button>
                </Stack>
                <LoginRedirect />
            </Stack>
        </form>
    );
}