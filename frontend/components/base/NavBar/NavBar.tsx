import {
    Box,
    Flex,
    useBreakpointValue,
} from '@chakra-ui/react';
import DesktopNav from './DesktopNav';
import MobileNav from './MobileNav';
import React from 'react';

// type User = {
//     username: string;
//     profilePic: string;
// }

// interface NavBarProps {
//     user: User | null;
//     isLoggedIn: boolean;
// }

export default function NavBar(props: any) {
    const { user, isLoggedIn } = props;
    const isDesktop = useBreakpointValue({ base: false, lg: true });
    const { username="test", profilePic="" } = user || {};
    return (
        <Box>
            <Flex py={2} px={4} borderBottom={1} align={'center'}>
                {isDesktop ? (
                    <DesktopNav profilePic={profilePic} isLoggedIn={isLoggedIn}/>
                ) : (
                    <MobileNav profilePic={profilePic} isLoggedIn={isLoggedIn}/>
                )}
            </Flex>
        </Box>
    );
}