import React from 'react';
import {
    Box,
    Button,
    Menu,
    MenuButton,
    MenuList,
    MenuItem,
} from '@chakra-ui/react';

import { BellIcon } from '@chakra-ui/icons';

export default function NotificationBell() {
    return (
        <Box>
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
                <MenuList
                >
                    <MenuItem>No new notifications</MenuItem>
                </MenuList>
            </Menu>
        </Box>
        
    )
}