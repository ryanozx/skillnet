import {
    Box,
    Flex,
    Text,
    Button,
    useColorModeValue,
    useBreakpointValue,
    HStack
} from '@chakra-ui/react';
import Searchbar from '../base/Searchbar';
import { useState } from 'react';
import ProfileButton from '../base/ProfileButton';
import NotificationBell from '../base/NotificationBell';


export default function DesktopNav () {  
    // const {isLoggedIn, setIsLoggedIn} = useState(true);
    const isLoggedIn = true

    return (
        
            <Flex flex={{ base: 1 }} justify={{ base: 'center', md: 'start' }} display={{base: 'none', md: 'flex'}} justifyContent={"space-between"}>
                <Text
                    pt={2}
                    textAlign={useBreakpointValue({ base: 'center', md: 'left' })}
                    fontFamily={'heading'}
                    color={useColorModeValue('gray.800', 'white')}>
                    FuckNet
                </Text>
                <Box w="50vw">
                    <Searchbar/>
                </Box>
                <HStack
                    spacing={6}
                >
                    {isLoggedIn ? 
                    
                    <>
                        <NotificationBell/>
                        <ProfileButton/>
                    </>
                    :
                    <>
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
                        <Button
                            as={'a'}
                            display={'inline-flex'}
                            fontSize={'sm'}
                            fontWeight={600}
                            color={'white'}
                            bg={'blue.400'}
                            href={'#'}
                            _hover={{
                                bg: 'blue.300',
                            }}
                        >
                            Sign Up
                        </Button>
                    </> 
                    }
                    
            </HStack>
            </Flex>
    );
};
