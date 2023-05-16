import React from 'react';
import {
    Avatar,
    Button,
    Menu,
    MenuButton,
    MenuList,
    MenuItem,
    MenuDivider,
} from '@chakra-ui/react';

import { BellIcon, SearchIcon } from '@chakra-ui/icons';

export default function NotificationBell() {
    return (
        <Menu>
            <MenuButton
                as={Button}
                rounded={'full'}
                variant={'link'}
                cursor={'pointer'}
                minW={0}
            >
                <BellIcon
                    boxSize={7}
                ></BellIcon>
            </MenuButton>
            <MenuList>
                <MenuItem>No new notifications</MenuItem>
            </MenuList>
        </Menu>
    )
}