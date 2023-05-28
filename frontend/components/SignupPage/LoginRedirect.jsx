import { Stack, Text, Link } from '@chakra-ui/react';

export default function LoginRedirect() {
    return (
        <Stack pt={6}>
            <Text align={'center'}>
                Already a user? <Link href='/login' color={'blue.400'}>Login</Link>
            </Text>
        </Stack>
    );
}