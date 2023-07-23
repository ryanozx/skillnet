import React from 'react';
import {
    Avatar,
    Box,
    Flex,
    HStack,
    VStack,
    Text,
    Heading,
} from '@chakra-ui/react';
import CropperComponent from '../base/EditProfilePicButton/CropperComponent';
import EditInfoBtn from './EditInfoBtn';
import { User } from '../../types';

interface BasicInfoProps {
    user: User;
    username: string;
    setUser?: React.Dispatch<React.SetStateAction<User>>;
}

export default function BasicInfo(props: BasicInfoProps) {
    return (
        <Box w="100%" px={10}>
            <Flex justifyContent={"space-between"} alignItems="flex-start">
                <HStack spacing={"10"}>
                {props.setUser ? (
                    <CropperComponent profilePic={props.user.ProfilePic} setUser={props.setUser} />
                    ) : (
                    <Avatar size="2xl" src={props.user.ProfilePic} />
                )}
                <VStack align="start">
                    <Heading size="md">{props.user.Name == "" ? "Anonymous User" : props.user.Name}</Heading>
                    <Text>{props.username}</Text>
                    <Text>{props.user.Title === "" ? "Title not available" : props.user.Title}</Text>
                </VStack>
                </HStack>
                <Flex alignSelf="flex-start">
                    {props.setUser && <EditInfoBtn user={props.user} setUser={props.setUser}/>}
                </Flex>
                
            </Flex>
        </Box>
    );
}