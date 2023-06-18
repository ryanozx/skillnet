import {
    Box,
    Flex,
    IconButton,
    useDisclosure,
    HStack,
    Drawer,
    DrawerBody,
    DrawerContent,
    DrawerCloseButton,
    DrawerOverlay,
    DrawerHeader,
} from '@chakra-ui/react';
import React from 'react';
import {
    HamburgerIcon,
    CloseIcon,

} from '@chakra-ui/icons';
import Searchbar from '../Searchbar';
import ProfileButton from './ProfileButton';
import NotificationBell from './NotificationBell';
import { useRef } from 'react'; 
import SideBar from '../SideBar/SideBar';


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
                <IconButton
                    onClick={onOpen}
                    icon={isOpen ? <CloseIcon w={3} h={3} /> : <HamburgerIcon w={5} h={5} />}
                    variant="ghost"
                    aria-label="Toggle Navigation"
                />
                <Box w="50%">
                    <Searchbar/>
                </Box>
                
                <HStack spacing={4}>

                    <NotificationBell/>
                    <ProfileButton profilePic = {profilePic} username={username}/>

                </HStack>
            </Flex>

            <Drawer
                isOpen={isOpen}
                placement='left'
                onClose={onClose}
                finalFocusRef={btnRef}
            >
                <DrawerOverlay />
                <DrawerContent>
                    <DrawerCloseButton />
                    <DrawerHeader>SkillNet</DrawerHeader>
                    <DrawerBody>
                        <SideBar/>
                    </DrawerBody>
                </DrawerContent>
            </Drawer>
        </>


    );
}
