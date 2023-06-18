import React from 'react';
import {
    Avatar,
    Button,
    MenuButton,
} from '@chakra-ui/react';

interface ProfileAvatarProps {
    profilePic: string;
}

const ProfileAvatar: React.FC<ProfileAvatarProps> = ({ profilePic }) => {
    return (
        <MenuButton
            as={Button}
            rounded={'full'}
            variant={'link'}
            cursor={'pointer'}
            minW={0}
        >
            <Avatar
                size={'md'}
                src={profilePic}
            />
        </MenuButton>
    )
}
export default ProfileAvatar;
