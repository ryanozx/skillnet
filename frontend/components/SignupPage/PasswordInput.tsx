import React, { useState , useEffect } from 'react';
import { 
    Button, 
    FormControl, 
    FormErrorMessage,
    FormLabel, 
    Input, 
    InputGroup, 
    InputRightElement } from '@chakra-ui/react';
import { ViewIcon, ViewOffIcon } from '@chakra-ui/icons';

interface PasswordInputProps {
    value: string;
    onChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
    passwordChanged: boolean;
    setPasswordChanged: React.Dispatch<React.SetStateAction<boolean>>;
    setPasswordError: React.Dispatch<React.SetStateAction<boolean>>;
}

export const PasswordInput: React.FC<PasswordInputProps> = (props : PasswordInputProps) => {
    const [showPassword, setShowPassword] = useState<boolean>(false);
    const isError = props.passwordChanged && !validatePassword(props.value);
    useEffect(() => props.setPasswordError(isError), [isError]);
    const onChange = (e : React.ChangeEvent<HTMLInputElement>) => {
        props.setPasswordChanged(true);
        props.onChange(e);
    }
    return (
        <FormControl id="password" isRequired isInvalid={isError}>
            <FormLabel>Password</FormLabel>
            <InputGroup>
                <Input 
                    data-testid="password-input"
                    type={showPassword ? 'text' : 'password'} 
                    name="password" 
                    value={ props.value } 
                    onChange={ onChange } />
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
            {isError && 
                <FormErrorMessage>Password must contain at least 8 characters, with at least one uppercase letter,
                    one lowercase letter, one number, and one special character.</FormErrorMessage>
            }
        </FormControl>
    )
}

function validatePassword(password: string) {
    const regex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/i;
    return regex.test(password)
}
