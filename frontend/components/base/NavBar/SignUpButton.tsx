import { Button } from '@chakra-ui/button';
import { useRouter } from 'next/router';
import React from 'react';

export default function SignUpButton() {
  const router = useRouter();

  const handleSignUp = () => {
    router.push('/signup');
  };

  return (
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
      onClick={handleSignUp}
    >
      Sign Up
    </Button>
  );
}
