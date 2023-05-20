import React from 'react';
import {
  Box,
  HStack,
  VStack,
  Text,
  Heading,
} from '@chakra-ui/react';
import EditPicBtn from './EditPicBtn';

export default function BasicInfo(user) {
    const { 
        name = "ivan tan", 
        username = "ivyy-poison", 
        description = "struggling cs student", 
        profilePic = "https://bit.ly/dan-abramov"
    } = user;

    return (
        <HStack spacing={"10"}>
            {/* <Avatar size="2xl" src={profilePic} /> */}
            <EditPicBtn currentProfilePic={profilePic}/>
            <VStack align="start">
                <Heading size="md">{name}</Heading>
                <Text>{username}</Text>
                <Text>{description}</Text>
            </VStack>
        </HStack>

    );
}