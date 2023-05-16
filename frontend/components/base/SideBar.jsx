import { Box, Flex, Icon, Link, Text, VStack } from '@chakra-ui/react';
import Searchbar from './Searchbar';

function Sidebar() {
  return (
    <Box
      position="fixed"
      top={20}
      left={0}
      h="100vh"
      w="15vw"
      bg="gray.200"
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
