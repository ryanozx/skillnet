import { ReactNode } from 'react';
import {
  Box,
  Flex,
  Avatar,
  HStack,
  Link,
  IconButton,
  Button,
  Menu,
  MenuButton,
  MenuList,
  MenuItem,
  MenuDivider,
  useDisclosure,
  useColorModeValue,
  Stack,
} from '@chakra-ui/react';
import { BellIcon, SearchIcon } from '@chakra-ui/icons';
import ProfileButton from './ProfileButton';
import NotificationBell from './NotificationBell';
import Searchbar from './Searchbar';

const Links = ['Dashboard', 'Projects', 'Team'];

const NavLink = (children) => (
  <Link
    px={2}
    py={1}
    rounded={'md'}
    _hover={{
      textDecoration: 'none',
      bg: 'gray.200'
    }}
    href={'#'}>
    {children}
  </Link>
);

export default function NavBar() {
    const { isOpen, onOpen, onClose } = useDisclosure();

    return (
        <>
            <Box bg={'gray.100'} px={4}>
                <Flex h={20} alignItems={'center'} justifyContent={'space-between'}>
                
                <HStack spacing={8} alignItems={'center'}>
                    <Box>SkillNet</Box>
                    <HStack
                    as={'nav'}
                    spacing={4}
                    display={{ base: 'none', md: 'flex' }}>
                    {/* {Links.map((link) => (
                        <NavLink key={link}>{link}</NavLink>
                    ))} */}
                    <Link
                        px={2}
                        py={1}
                        rounded={'md'}
                        _hover={{
                        textDecoration: 'none',
                        bg: 'gray.200'
                        }}
                        href={'#'}
                    > t1 </Link>                    
                </HStack>
                
                <Searchbar/>

                </HStack>
                    <Flex alignItems={'center'}>
                        <HStack
                            spacing={7}    
                        >
                            <NotificationBell></NotificationBell>
                            <ProfileButton></ProfileButton>
                        </HStack>
                    </Flex>
                </Flex>

                
            </Box>
        </>
  );
}