import React from 'react';
import {
  Box,
  VStack,
  Text,
  Heading,
} from '@chakra-ui/react';

export default function AboutMe(user) {
    const {
        about = "No description available"
    } = user
    return (
        <Box
            w="100%"
            outline={"4px solid gray"}
            p={4}
            h={"20vh"}
            maxH={"20vh"}
        >
            <VStack spacing={5} align="start">
                <Heading size="md">About Me</Heading>
                <Box 
                    bg="green.200"
                    w="100%"
                    px={4}
                >
                    <Text>{about}</Text>
                </Box>
            </VStack>
        </Box>
    );
}