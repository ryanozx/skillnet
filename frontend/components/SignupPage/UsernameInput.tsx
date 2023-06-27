import React, { ChangeEvent } from 'react';
import { 
    FormControl, 
    FormLabel, 
    Input
} from '@chakra-ui/react';


interface UsernameInputProps {
    value: string;
    onChange: (newUsername: string) => void;
}

export const UsernameInput: React.FC<UsernameInputProps> = ({ value, onChange }) => (
    <FormControl id="username" isRequired>
        <FormLabel>Username</FormLabel>
        <Input type="text" value={value} onChange={e => onChange(e.target.value)} />
    </FormControl>
);

