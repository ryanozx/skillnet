import { Button } from '@chakra-ui/button';

export default function LogInButton() {
    return (
        <Button
            as={'a'}
            fontSize={'sm'}
            fontWeight={400}
            color={'blackAlpha.900'}
            variant={'link'}
            href={'#'}
        >
            Log in
        </Button>
    );
}