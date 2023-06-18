import { FormControl, 
    FormLabel, 
    Input 
} from '@chakra-ui/react';
import React from "react";

interface UsernameInputProps {
    value: string;
    onChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
}

export const UsernameInput: React.FC<UsernameInputProps> = ({ value, onChange }) => {
    return (
        <FormControl id="username">
            <FormLabel>username</FormLabel>
            <Input
                type="username"
                name="username"
                value={value}
                onChange={onChange}
            />
        </FormControl>
    );
}
