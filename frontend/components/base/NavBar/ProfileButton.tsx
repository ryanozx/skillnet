import React, { MouseEventHandler } from 'react';
import {
    Link,
    Menu,
    MenuList,
    MenuItem,
    MenuDivider,
} from '@chakra-ui/react';
import LogoutButton from './LogoutButton';
import ProfileAvatar from './ProfileAvatar';

export interface ProfileButtonProps {
    profilePic: string;
    username: string;
}

export default function ProfileButton(props: ProfileButtonProps) {
    const { profilePic, username } = props;
    return (
        <Menu>
            <ProfileAvatar profilePic={profilePic}/>
            <MenuList>
            <Link href={`/profile/${username}`}><MenuItem>View your profile</MenuItem></Link>
                <MenuDivider />
                <LogoutButton />
            </MenuList>
        </Menu>
    );
}
