import React, { ChangeEvent } from 'react';
import { 
    FormControl, 
    FormLabel, 
    Input
} from '@chakra-ui/react';


interface UsernameInputProps {
    value: string;
    onChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
}

export const UsernameInput: React.FC<UsernameInputProps> = ({ value, onChange }) => (
    <FormControl id="username" isRequired>
        <FormLabel>Username</FormLabel>
        <Input data-testid="username-input" type="text" name="username" value={value} onChange={onChange} />
    </FormControl>
);

