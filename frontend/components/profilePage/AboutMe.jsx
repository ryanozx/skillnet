import React from 'react';
import {
  Box,
  HStack,
  VStack,
  Text,
  Heading,
} from '@chakra-ui/react';
import EditPicBtn from './EditPicBtn';

export default function AboutMe(user) {
    const {about} = user
    return (
        <Box
            w="100%"
            outline={"4px solid gray"}
            p={4}
            h={"20vh"}
            maxH={"20vh"}
        >
            <VStack spacing={5} align="start">
                <Heading size="md">About</Heading>
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