import React, { useState, MouseEventHandler, ChangeEvent } from 'react';
import { Stack, Button, FormControl, FormLabel, Input, InputGroup, InputRightElement, useToast } from '@chakra-ui/react';
import LoginRedirect from './LoginRedirect';
import { ViewIcon, ViewOffIcon } from '@chakra-ui/icons';
import axios from 'axios';
import { useRouter } from 'next/router';

interface EmailInputProps {
    value: string;
    onChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
}

export const EmailInput: React.FC<EmailInputProps> = ({ value, onChange }) => (
    <FormControl id="email" isRequired>
        <FormLabel>Email address</FormLabel>
        <Input data-testid="email-input" type="email" name="email" value={value} onChange={onChange} />
    </FormControl>
);

