import {
    Box,
    Flex,
    useBreakpointValue,
} from '@chakra-ui/react';
import DesktopNav from './DesktopNav';
import MobileNav from './MobileNav';
  
export default function NavBar( {user, isLoggedIn} ) {
    const isDesktop = useBreakpointValue({ base: false, lg: true });
    const {username="test", profilePic=""} = user;
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