import { Box, Flex, Stack } from '@chakra-ui/react';
import FormHeading from './FormHeading';
import SignUpForm from './SignUpForm';

export default function SignUpPageContainer() {
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
}