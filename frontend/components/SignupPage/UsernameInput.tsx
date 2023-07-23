import React, {useEffect} from 'react';
import { 
    FormControl, 
    FormErrorMessage,
    FormLabel, 
    Input
} from '@chakra-ui/react';


interface UsernameInputProps {
    value: string;
    onChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
    usernameChanged: boolean;
    setUsernameChanged: React.Dispatch<React.SetStateAction<boolean>>;
    setUsernameError: React.Dispatch<React.SetStateAction<boolean>>;
}

export const UsernameInput: React.FC<UsernameInputProps> = (props : UsernameInputProps) => {
    const isEmpty = props.usernameChanged && props.value === "";
    useEffect(() => props.setUsernameError(isEmpty), [isEmpty]);

    const onChange = (e : React.ChangeEvent<HTMLInputElement>) => {
        props.setUsernameChanged(true);
        props.onChange(e);
    }
    return (
        <FormControl id="username" isRequired isInvalid={isEmpty}>
            <FormLabel>Username</FormLabel>
            <Input 
                data-testid="username-input" 
                type="text" 
                name="username" 
                value={props.value} 
                onChange={onChange} />
            {isEmpty && <FormErrorMessage>Username cannot be empty</FormErrorMessage>}
        </FormControl>)
};

