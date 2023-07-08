import { 
    FormControl, 
    FormLabel, 
    Input, 
} from '@chakra-ui/react';
import React from "react";

interface PasswordInputProps {
    value: string;
    onChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
}

export const PasswordInput: React.FC<PasswordInputProps> = ({ value, onChange }) => {
    return (
        <FormControl id="password">
            <FormLabel>Password</FormLabel>
            <Input
                data-testid="password-input"
                type="password"
                name="password"
                value={value}
                onChange={onChange}
            />
        </FormControl>
    );
}
