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
    Input,
} from '@chakra-ui/react';

import {
    HamburgerIcon,
    CloseIcon,

} from '@chakra-ui/icons';
import Searchbar from '../base/Searchbar';
import ProfileButton from './ProfileButton';
import NotificationBell from './NotificationBell';
import LogInButton from './LogInButton';
import SignUpButton from './SignUpButton';
import { useRef } from 'react'; 
import SideBar from '../SideBar/SideBar';


export default function MobileNav(props) {
    const { isOpen, onOpen, onClose } = useDisclosure()
    const btnRef = useRef()
    const { isLoggedIn=true, profilePic='' } = props;
    console.log(isLoggedIn + " is logged in?");
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
                
                <HStack spacing={isLoggedIn ? 4 : 2}>

                    <NotificationBell/>
                    <ProfileButton profilePic = {profilePic}/>

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
