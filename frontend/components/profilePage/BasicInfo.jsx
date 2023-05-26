import React from 'react';
import {
  Box,
  Flex,
  HStack,
  VStack,
  Text,
  Heading,
} from '@chakra-ui/react';
import EditPicBtn from './EditPicBtn';
import EditInfoBtn from './EditInfoBtn';

export default function BasicInfo(user) {
    const { 
        name = "ivan tan", 
        username = "ivyy-poison", 
        description = "struggling cs student", 
        profilePic = "https://bit.ly/dan-abramov"
    } = user;

    return (
        <Box w="100%" px={10}>
            <Flex justifyContent={"space-between"} alignItems="flex-start">
                <HStack spacing={"10"}>
                {/* <Avatar size="2xl" src={profilePic} /> */}
                <EditPicBtn currentProfilePic={profilePic}/>
                    <VStack align="start">
                        <Heading size="md">{name}</Heading>
                        <Text>{username}</Text>
                        <Text>{description}</Text>
                    </VStack>
                </HStack>
                <Flex alignSelf="flex-start">
                    <EditInfoBtn/>
                </Flex>
                
            </Flex>
        </Box>
    );
}