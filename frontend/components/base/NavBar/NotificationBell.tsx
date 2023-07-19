import React from 'react';
import { Box, Button, Menu, MenuButton, MenuItem, MenuList } from '@chakra-ui/react';
import { BellIcon } from '@chakra-ui/icons';

interface NotificationBellProps {
    notifications: string[];
    hasNewNotifications: boolean;
    setHasNewNotifications: React.Dispatch<React.SetStateAction<boolean>>;
}

export default function NotificationBell(props: NotificationBellProps) {
    const { notifications = [], hasNewNotifications, setHasNewNotifications } = props;

    return (
        <Box>
            <Menu>
                <MenuButton
                    as={Button}
                    rounded={'full'}
                    variant={'link'}
                    cursor={'pointer'}
                    minW={0}
                    onClick={() => {setHasNewNotifications(false)}}
                >
                    <BellIcon
                        boxSize={7}
                        color={hasNewNotifications ? "red.500" : "gray.500"} // Change color if there are new notifications
                    />
                </MenuButton>
                <MenuList
                    maxH="400px" // Limit the height to show a maximum of 10 items (assuming each item is 40px high)
                    overflowY="scroll" // Enable vertical scrolling
                    maxW="300px"
                >
                    {notifications.length > 0
                        ? notifications.map((notification, index) => (
                            <MenuItem key={index}>{notification}</MenuItem>
                        ))
                        : <MenuItem>No new notifications</MenuItem>
                    }
                </MenuList>
            </Menu>
        </Box>  
    )
}
