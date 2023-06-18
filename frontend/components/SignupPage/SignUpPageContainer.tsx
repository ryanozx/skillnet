import React from 'react';
import { Box, Flex, Stack } from '@chakra-ui/react';
import FormHeading from './FormHeading';
import SignUpForm from './SignUpForm';
import { preventAuthAccess } from '../../WithAuthRedirect';


export default preventAuthAccess(function SignUpPageContainer() {
    return (
        <Flex
            minH={'100vh'}
            align={'center'}
            justify={'center'}
            bg='green.300'>
            <Stack spacing={8} mx={'auto'} py={12} px={6}>
                <FormHeading />
                <Box
                    rounded={'lg'}
                    bg='gray.100'
                    boxShadow={'lg'}
                    p={8}
                    w={{ base: '90vw', md: '60vw', lg: '30vw' }}
                >
                    <SignUpForm />
                </Box>
            </Stack>
        </Flex>
    );
});