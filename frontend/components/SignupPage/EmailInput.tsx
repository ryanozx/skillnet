import React, {useEffect} from 'react';
import {
    FormControl, 
    FormErrorMessage,
    FormLabel, 
    Input} from '@chakra-ui/react';

interface EmailInputProps {
    value: string;
    onChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
    emailChanged: boolean;
    setEmailChanged: React.Dispatch<React.SetStateAction<boolean>>;
    setEmailError: React.Dispatch<React.SetStateAction<boolean>>;
}

export const EmailInput: React.FC<EmailInputProps> = (props : EmailInputProps) => {
    const isError = props.emailChanged && !validateEmail(props.value);
    useEffect(() => props.setEmailError(isError), [isError]);
    const onChange = (e : React.ChangeEvent<HTMLInputElement>) => {
        props.setEmailChanged(true);
        props.onChange(e);
    }
    return (<FormControl id="email" isRequired isInvalid={isError}>
        <FormLabel>Email address</FormLabel>
        <Input data-testid="email-input" type="email" name="email" value={props.value} onChange={onChange} />
        {isError && 
            <FormErrorMessage>Please enter a valid email address.</FormErrorMessage>
        }
    </FormControl>);
};

function validateEmail(email: string) {
    const regex = /^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/i;
    return regex.test(email)
}