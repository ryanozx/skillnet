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
import ProfileButton from './ProfileButton';
import NotificationBell from './NotificationBell';
import LogInButton from './LogInButton';
import SignUpButton from './SignUpButton';


export default function DesktopNav () {  
    // const {isLoggedIn, setIsLoggedIn} = useState(true);
    const isLoggedIn = true

    return (
        
            <Flex flex={{ base: 1 }} justify={{ base: 'center', md: 'start' }} display={{base: 'none', md: 'flex'}} justifyContent={"space-between"}>
                <Text
                    pt={2}
                    textAlign={useBreakpointValue({ base: 'center', md: 'left' })}
                    fontFamily={'heading'}
                    color={useColorModeValue('gray.800', 'white')}    
                >
                    SkillNet
                </Text>
                <Box w="50vw">
                    <Searchbar/>
                </Box>
                <HStack spacing={6}>
                    {isLoggedIn ?                
                    <>
                        <NotificationBell/>
                        <ProfileButton/>
                    </>
                    :
                    <>
                        <LogInButton/>
                        <SignUpButton/>
                    </> 
                    } 
                </HStack>
            </Flex>
    );
};
