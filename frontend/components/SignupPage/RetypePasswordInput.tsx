import React, { useState , useEffect } from 'react';
import { 
    Button, 
    FormControl, 
    FormErrorMessage,
    FormLabel, 
    Input, 
    InputGroup, 
    InputRightElement } from '@chakra-ui/react';
import { 
    ViewIcon, 
    ViewOffIcon } from '@chakra-ui/icons';

interface RetypePasswordInputProps {
    value: string;
    onChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
    passwordMismatch: boolean;
    retypePasswordChanged: boolean;
    setRetypePasswordChanged: React.Dispatch<React.SetStateAction<boolean>>;
    setRetypePasswordError: React.Dispatch<React.SetStateAction<boolean>>;
}

export const RetypePasswordInput: React.FC<RetypePasswordInputProps> = (props : RetypePasswordInputProps) => {
    const [showPassword, setShowPassword] = useState<boolean>(false);
    const isError = props.retypePasswordChanged && props.passwordMismatch;
    useEffect(() => props.setRetypePasswordError(isError), [isError]);
    const onChange = (e : React.ChangeEvent<HTMLInputElement>) => {
        props.setRetypePasswordChanged(true);
        props.onChange(e);
    }
    return (
        <FormControl id="passwordRetype" isRequired isInvalid={isError}>
            <FormLabel>Retype Password</FormLabel>
            <InputGroup>
                <Input 
                    data-testid="retypePassword-input"
                    type={showPassword ? 'text' : 'password'} 
                    name="retypePassword" 
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
            {isError && <FormErrorMessage>Passwords do not match</FormErrorMessage>}
        </FormControl>
    )
}

