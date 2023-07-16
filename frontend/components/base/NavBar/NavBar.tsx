import {
    Box,
    Flex,
    useBreakpointValue,
} from '@chakra-ui/react';
import DesktopNav from './DesktopNav';
import MobileNav from './MobileNav';
import React from 'react';

interface NavBarProps {
    profilePic: string;
    username: string;
}

export default function NavBar(props: NavBarProps) {
    const { profilePic, username } = props;
    const isDesktop = useBreakpointValue({ base: false, lg: true });
    
    
    return (
        <Box>
            <Flex py={2} px={4} borderBottom={1} align={'center'}>
                {isDesktop ? (
                    <DesktopNav profilePic={profilePic} username={username}/>
                ) : (
                    <MobileNav profilePic={profilePic} username={username}/>
                )}
            </Flex>
        </Box>
    );
}