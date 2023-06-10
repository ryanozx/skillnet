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
import React, { useState } from 'react';
import ProfileButton from './ProfileButton';
import NotificationBell from './NotificationBell';

interface DesktopNavProps {
    profilePic: string;
    isLoggedIn: boolean;
}

export default function DesktopNav (props: DesktopNavProps) {  
    const { isLoggedIn = true, profilePic = '' } = props;
    console.log(isLoggedIn + " is logged in?");
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
            <HStack spacing={isLoggedIn? 6 : 3}>
                <NotificationBell/>
                <ProfileButton profilePic={profilePic}/>

            </HStack>
        </Flex>
    );
};
