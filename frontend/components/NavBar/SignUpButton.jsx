import Link from 'next/link';
import { Button } from '@chakra-ui/button';

export default function SignUpButton() {
    return (
        <Link href="/signup">
            <Button
                as={'a'}
                display={'inline-flex'}
                fontSize={'sm'}
                fontWeight={600}
                color={'white'}
                bg={'blue.400'}
                _hover={{
                    bg: 'blue.300',
                }}
            >
                Sign Up
            </Button>
        </Link>
    );
}
