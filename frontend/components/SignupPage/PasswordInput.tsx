import React, { useState, ChangeEvent } from 'react';
import { Stack, Button, FormControl, FormLabel, Input, InputGroup, InputRightElement, useToast } from '@chakra-ui/react';
import LoginRedirect from './LoginRedirect';
import { ViewIcon, ViewOffIcon } from '@chakra-ui/icons';
import axios from 'axios';
import { useRouter } from 'next/router';

interface PasswordInputProps {
    value: string;
    onChange: (newPassword: string) => void;
}

export const PasswordInput: React.FC<PasswordInputProps> = ({ value, onChange }) => {
    const [showPassword, setShowPassword] = useState<boolean>(false);

    return (
        <FormControl id="password" isRequired>
            <FormLabel>Password</FormLabel>
            <InputGroup>
                <Input aria-label="Password" type={showPassword ? 'text' : 'password'} value={value} onChange={e => onChange(e.target.value)} />
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
    )
}

