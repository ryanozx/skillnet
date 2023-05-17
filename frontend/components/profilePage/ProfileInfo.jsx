import React from 'react';
import {
  Box,
  HStack,
  VStack,
  Avatar,
  Text,
  Heading,
  IconButton,
  Flex
} from '@chakra-ui/react';
import { EditIcon } from '@chakra-ui/icons';

export default function ProfileInfo(user) {
    const { 
        name = "ivan tan", 
        username = "ivyy-poison", 
        description = "struggling cs student", 
        about = "asomsef", 
        profilePic = "https://bit.ly/dan-abramov"
    } = user;

  return (
    <Box p={4}>
        <VStack spacing={10} align="start">

            <HStack spacing={"10"}>
                <Avatar size="2xl" src={profilePic} />
                <VStack align="start">
                    <Heading size="md">{name}</Heading>
                    <Text>{username}</Text>
                    <Text>{description}</Text>
                </VStack>
            </HStack>

            <Box
                w="100%"
                outline={"1px solid black"}
                p={4}
                h={"20vh"}
                maxH={"20vh"}
            >
                <VStack spacing={5} align="start">
                    <Heading size="md">About</Heading>
                    <Box 
                        bg="gray.100"
                        w="100%"
                        px={4}
                    >
                        <Text>{about}</Text>
                    </Box>
                    
                </VStack>
                
            </Box>

            <IconButton
                alignSelf="flex-end"
                icon={<EditIcon />}
                aria-label="Edit profile"
                variant="outline"
            />
        </VStack>

    </Box>
    
  );
};
