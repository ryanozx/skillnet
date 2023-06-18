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
    name?: string;
    username: string;
    title?: string;
    profilePic?: string;
    setUser?: React.Dispatch<React.SetStateAction<User>>;
}

export default function BasicInfo(props: BasicInfoProps) {
    const { name, username, title, profilePic, setUser } = props;
    return (
        <Box w="100%" px={10}>
            <Flex justifyContent={"space-between"} alignItems="flex-start">
                <HStack spacing={"10"}>
                {setUser ? (
                    <CropperComponent profilePic={profilePic} setUser={setUser} />
                    ) : (
                    <Avatar size="2xl" src={profilePic} />
                )}
                <VStack align="start">
                    <Heading size="md">{name}</Heading>
                    <Text>{username}</Text>
                    <Text>{title}</Text>
                </VStack>
                </HStack>
                <Flex alignSelf="flex-start">
                    {setUser && <EditInfoBtn setUser={setUser}/>}
                </Flex>
                
            </Flex>
        </Box>
    );
}