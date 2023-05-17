import { Box, Flex, Icon, Link, Text, VStack } from '@chakra-ui/react';
import Searchbar from './Searchbar';

function Sidebar() {
  return (
    <Box

      h="100vh"
      w="100%"
      color="black"
      display="flex"
      flexDirection="column"
      alignItems="center"
      justifyContent="space-between"
      py={4}
      px={2}
      zIndex={999}
    >
      <VStack  spacing={4} alignItems="center">
        <Searchbar></Searchbar>
        <Link href="/">
          link1
        </Link>
      </VStack>
      <Flex alignItems="center" justifyContent="center">
        <Text fontSize="sm">Your App</Text>
      </Flex>
    </Box>
  );
}

export default Sidebar;
