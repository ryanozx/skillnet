import {
    Box,
    Flex,
    useDisclosure,
    HStack,
} from '@chakra-ui/react';
import React from 'react';
import Searchbar from '../Searchbar';
import ProfileButton from './ProfileButton';
import NotificationBell from './NotificationBell';
import { useRef } from 'react'; 
import MobileHamburgerButton from './MobileHamburgerButton';
import MobileDrawerMenu from './MobileDrawerMenu';


interface MobileNavProps {
    profilePic: string;
    username: string;
    notifications: string[];
    hasNewNotifications: boolean;
    setHasNewNotifications: React.Dispatch<React.SetStateAction<boolean>>;
}   

export default function MobileNav(props: MobileNavProps) {
    const { isOpen, onOpen, onClose } = useDisclosure()
    const btnRef = useRef(null)
    const { profilePic, username, notifications, hasNewNotifications, setHasNewNotifications } = props;
    return (
        <>
            <Flex flex={1} display='flex' justifyContent={"space-between"}>
                <MobileHamburgerButton onOpen={onOpen} />
                <Box w="50%">
                    <Searchbar/>
                </Box>
                <HStack spacing={4}>
                    <NotificationBell 
                    notifications={notifications} 
                    hasNewNotifications={hasNewNotifications} 
                    setHasNewNotifications={setHasNewNotifications}/>
                    <ProfileButton profilePic = {profilePic} username={username}/>
                </HStack>
            </Flex>
            <MobileDrawerMenu isOpen={isOpen} onClose={onClose} btnRef={btnRef} />
        </>
    );
}
