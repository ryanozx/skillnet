import {
    Box,
    Flex,
    useBreakpointValue,
    useDisclosure,
} from '@chakra-ui/react';


import { useState } from 'react';
import DesktopNav from './DesktopNav';
import MobileNav from './MobileNav';
  
export default function NavBar() {
    const { isOpen, onToggle } = useDisclosure();
    const { isLoggedIn, setIsLoggedIn } = useState(false);
    const isDesktop = useBreakpointValue({ base: false, lg: true });
  
    return (
        <Box>
            <Flex
                py={2}
                px={4}
                borderBottom={1}
                align={'center'}
            >
                {isDesktop ? <DesktopNav /> : <MobileNav />}
            </Flex>
        </Box>
    );
  }
  

  