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

export default function NavBar() {
    const { isOpen, onOpen, onClose } = useDisclosure();

    return (
        <>
            <Box px={4}>
                <Flex h="7vh" alignItems={'center'} justifyContent={'space-around'}>
                
                    <HStack spacing={8} alignItems={'center'}>
                        <Box>SkillNet</Box>
                    </HStack>
                
                    <Searchbar/>

                    <HStack>
                        <Flex alignItems={'center'}>
                            <HStack
                                spacing={7}    
                            >
                                <NotificationBell></NotificationBell>
                                <ProfileButton></ProfileButton>
                            </HStack>
                        </Flex>
                    </HStack>
                </Flex>


                
            </Box>
        </>
  );
}