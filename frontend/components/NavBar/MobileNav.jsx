import {
    Box,
    Flex,
    IconButton,
    Button,
    useDisclosure,
    HStack
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


export default function MobileNav() {
  const { isOpen, onToggle } = useDisclosure();
  const isLoggedIn = true

  return (
    <Flex flex={1} display={{ base: 'flex', md: 'none' }} justifyContent={"space-between"}>
        <IconButton
            onClick={onToggle}
            icon={isOpen ? <CloseIcon w={3} h={3} /> : <HamburgerIcon w={5} h={5} />}
            variant="ghost"
            aria-label="Toggle Navigation"
        />
        <Box>
            <Searchbar/>
        </Box>
        
        <HStack spacing={4}>
            {isLoggedIn ?  
                <>
                    <NotificationBell/>
                    <ProfileButton/>
                </>
                :
                <>
                    <LogInButton/>
                    <SignInButton/>    
                </> 
            }
        </HStack>
    </Flex>
  );
}
