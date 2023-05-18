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
    const isDesktop = useBreakpointValue({ base: false, md: true });
  
    return (
        <Box>
            <Flex
                minH={'60px'}
                py={{ base: 2 }}
                px={{ base: 4 }}
                borderBottom={1}
                align={'center'}
            >
                {isDesktop ? <DesktopNav /> : <MobileNav />}
            </Flex>
        </Box>
    );
  }
  

  