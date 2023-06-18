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
}   

export default function MobileNav(props: MobileNavProps) {
    const { isOpen, onOpen, onClose } = useDisclosure()
    const btnRef = useRef(null)
    const { profilePic, username } = props;
    return (
        <>
            <Flex flex={1} display='flex' justifyContent={"space-between"}>
                <MobileHamburgerButton onOpen={onOpen} />
                <Box w="50%">
                    <Searchbar/>
                </Box>
                <HStack spacing={4}>
                    <NotificationBell/>
                    <ProfileButton profilePic = {profilePic} username={username}/>
                </HStack>
            </Flex>
            <MobileDrawerMenu isOpen={isOpen} onClose={onClose} btnRef={btnRef} />
        </>
    );
}
