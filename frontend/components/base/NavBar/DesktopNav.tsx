import {
    Box,
    Flex,
    Link,
    Text,
    HStack
} from '@chakra-ui/react';
import Searchbar from '../Searchbar';
import React, { useState } from 'react';
import ProfileButton from './ProfileButton';
import NotificationBell from './NotificationBell';

interface DesktopNavProps {
    profilePic: string;
    username: string;
}

export default function DesktopNav (props: DesktopNavProps) {  
    const { profilePic, username } = props;
    return (
        <Flex 
            flex={1}  
            justify={{ base: 'center', md: 'start' }} 
            display={{base: 'none', md: 'flex'}} 
            justifyContent={"space-between"}>
            <Link href="/feed"> 
                <Text
                    pt={2}
                    textAlign={{ base: 'center', md: 'left' }}
                    fontFamily='heading'
                    color='gray.800'>
                    SKILLNET
                </Text>
            </Link>
            <Box w="50vw">
                <Searchbar/>
            </Box>
            <HStack spacing={6}>
                <NotificationBell/>
                <ProfileButton profilePic={profilePic} username={username}/>
            </HStack>
        </Flex>
    );
};
