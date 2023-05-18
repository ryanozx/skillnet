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
import ProfileButton from '../base/ProfileButton';
import NotificationBell from '../base/NotificationBell';



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
                    <Button
                        as={'a'}
                        fontSize={'sm'}
                        fontWeight={400}
                        color={'blackAlpha.900'}
                        variant={'link'}
                        href={'#'}
                    >
                        Log in
                    </Button>
                    <Button
                        as={'a'}
                        display={'inline-flex'}
                        fontSize={'sm'}
                        fontWeight={600}
                        color={'white'}
                        bg={'blue.400'}
                        href={'#'}
                        _hover={{
                            bg: 'blue.300',
                        }}
                    >
                        Sign Up
                    </Button>
                </> 
            }
        </HStack>
    </Flex>
  );
}
