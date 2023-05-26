import { Button } from '@chakra-ui/react';
import { useRouter } from 'next/router';

export default function LogInButton() {
    const router = useRouter();

    const handleLogin = () => {
        router.push('/login');
    };

    return (
        <Button
            fontSize={'sm'}
            fontWeight={400}
            color={'blackAlpha.900'}
            variant={'link'}
            onClick={handleLogin}
        >
        Log in
        </Button>
    );
}

