import React, { MouseEventHandler } from 'react';
import {Stack, FormControl, FormLabel, Input, Button, Checkbox} from '@chakra-ui/react';
import SignUpSection from './SignUpSection';

interface ActionSectionProps {
    onSubmit: React.MouseEventHandler<HTMLButtonElement>;
}

export default function ActionSection({onSubmit}: ActionSectionProps) {
    return (
        <Stack spacing={10}>
            <Checkbox>Remember me</Checkbox>
            <Button
                bg={'blue.400'}
                color={'white'}
                _hover={{
                    bg: 'blue.500',
                }}
                onClick={onSubmit}
            >
                Sign in
            </Button>
            <SignUpSection />
        </Stack>
    );
}