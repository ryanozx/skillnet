import {
    Box,
    Flex,
    useBreakpointValue,
} from '@chakra-ui/react';
import DesktopNav from './DesktopNav';
import MobileNav from './MobileNav';
  
function NavBar( {user, isLoggedIn} ) {
    const isDesktop = useBreakpointValue({ base: false, lg: true });
    const {username="test", profilePic=""} = user;
    console.log("we at the navbar here, the username should not be test if logged in " + username)
    console.log("at navbar, isLoggedIn currently should be False " + isLoggedIn)
    return (
        <Box>
            <Flex
                py={2}
                px={4}
                borderBottom={1}
                align={'center'}
            >
                {isDesktop ? <DesktopNav profilePic={profilePic}/> : <MobileNav profilePic={profilePic}/>}
            </Flex>
        </Box>
    );
}

// const mapStateToProps = state => ({
//     user: state.user,
// });
  

// export default connect(mapStateToProps)(NavBar);
export default NavBar;