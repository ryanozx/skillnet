import LoginPageContainer from '../../components/LoginPage/LoginPageContainer';
  
export default function LoginPage() {
    return (
        // <Flex
        //     minH={'100vh'}
        //     align={'center'}
        //     justify={'center'}
        //     bg={useColorModeValue('gray.50', 'gray.800')}>
        //     <Stack spacing={8} mx={'auto'} py={12} px={6}>
        //     <Stack align={'center'}>
        //         <Heading fontSize={'4xl'}>Sign in to your account</Heading>
        //         <Text fontSize={'lg'} color={'gray.600'}>
        //             to enjoy all of our cool features ✌️
        //         </Text>
        //     </Stack>
        //     <Box
        //         rounded={'lg'}
        //         bg={useColorModeValue('white', 'gray.700')}
        //         boxShadow={'lg'}
        //         p={8}
        //         w = {{base:'90vw', md:'60vw', lg:'30vw'}}
        //         >
        //         <Stack spacing={4}>
        //             <FormControl id="email">
        //                 <FormLabel>Email address</FormLabel>
        //                 <Input type="email" />
        //             </FormControl>
        //             <FormControl id="password">
        //                 <FormLabel>Password</FormLabel>
        //                 <Input type="password" />
        //             </FormControl>
        //             <Stack spacing={10}>
        //                     <Checkbox>Remember me</Checkbox>
        //                 <Button
        //                     bg={'blue.400'}
        //                     color={'white'}
        //                     _hover={{
        //                         bg: 'blue.500',
        //                 }}>
        //                     Sign in
        //                 </Button>
        //             </Stack>
        //             <Stack pt={6}>
        //                 <Text align={'center'}>
        //                     Don't have an account? <Link href='/signup' color={'blue.400'}>Sign up</Link>
        //                 </Text>
        //             </Stack>
        //         </Stack>
        //     </Box>
        //     </Stack>
        // </Flex>
        <LoginPageContainer/>
    );
}