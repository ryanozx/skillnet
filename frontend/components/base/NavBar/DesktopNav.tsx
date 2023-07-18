import {
    Box,
    Flex,
    Link,
    Text,
    HStack
} from '@chakra-ui/react';
import Searchbar from '../Searchbar';
import React from 'react';
import ProfileButton from './ProfileButton';
import NotificationBell from './NotificationBell';

interface DesktopNavProps {
    profilePic: string;
    username: string;
    notifications: string[];
    hasNewNotifications: boolean;
    setHasNewNotifications: React.Dispatch<React.SetStateAction<boolean>>;
}

export default function DesktopNav (props: DesktopNavProps) {  
    const { profilePic, username, notifications, hasNewNotifications, setHasNewNotifications } = props;
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
                <NotificationBell 
                    notifications={notifications} 
                    hasNewNotifications={hasNewNotifications} 
                    setHasNewNotifications={setHasNewNotifications}/>
                <ProfileButton profilePic={profilePic} username={username}/>
            </HStack>
        </Flex>
    );
};
