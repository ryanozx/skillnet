import Link from 'next/link';
import { Button } from '@chakra-ui/react';

export default function LogInButton() {
    return (
        <Link href="/login">
            <Button
                fontSize={'sm'}
                fontWeight={400}
                color={'blackAlpha.900'}
                variant={'link'}
                as={'a'}
            >
                Log in
            </Button>
        </Link>
    );
}
