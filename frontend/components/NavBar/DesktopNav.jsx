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


export default function DesktopNav (props) {  
    const { isLoggedIn = true, profilePic = '' } = props;
    return (
        <Flex 
            flex={1}  
            justify={{ base: 'center', md: 'start' }} 
            display={{base: 'none', md: 'flex'}} 
            justifyContent={"space-between"}>
            <Text
                pt={2}
                textAlign={{ base: 'center', md: 'left' }}
                fontFamily='heading'
                color='gray.800'>
                SKILLNET
            </Text>
            <Box w="50vw">
                <Searchbar/>
            </Box>
            <HStack spacing={6}>
                {isLoggedIn ?  (
                    <>
                    <NotificationBell/>
                    <ProfileButton profilePic={profilePic}/>
                    </>
                ) : (
                    <>
                    <LogInButton/>
                    <SignUpButton/>
                    </> 
                )}
            </HStack>
        </Flex>
    );
};
