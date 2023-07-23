import React, { MouseEventHandler } from 'react';
import { Stack, Button, useToast } from '@chakra-ui/react';


interface SubmitButtonProps {
    formValid: boolean;
}

export const SignUpButton = (props : SubmitButtonProps) => {
    return(
        <Stack spacing={10} pt={2}>
            <Button
                data-testid="signup-button"
                isDisabled={!props.formValid}
                type="submit"
                size="lg"
                bg={'blue.400'}
                color={'white'}
                _hover={{
                    bg: 'blue.500',
                }}>
                Sign up
            </Button>
        </Stack>
    )
};
