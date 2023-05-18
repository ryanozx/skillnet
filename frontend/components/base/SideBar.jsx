import { Box, Flex, Icon, Link, Text, VStack } from '@chakra-ui/react';
import Searchbar from './Searchbar';

export default function Sidebar() {
  return (
    <Box
        w="100%"
        color="black"
        alignItems="start"
        p={5}
        zIndex={9}
    >
        <VStack>
            <Searchbar></Searchbar>
            <Link href="/">
                link1
            </Link>
        </VStack>
    </Box>
  );
}
