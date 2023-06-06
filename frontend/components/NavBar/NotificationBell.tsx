import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Box, Button, Menu, MenuButton, MenuItem, MenuList } from '@chakra-ui/react';
import { BellIcon } from '@chakra-ui/icons';

export default function NotificationBell() {
    const [notifications, setNotifications] = useState([]);
    const url = '/fake-url';
    console.log('API call to get notifications for user');
    useEffect(() => {
        const sessionId = sessionStorage.getItem('sessionId');
        axios.post(url, {
            headers: {
                'Authorization': `Bearer ${sessionId}`
            }
        })
        .then(result => {
            setNotifications(result.data);
        })
        .catch(error => {
            // console.error(error);
            setNotifications([]);
        });

    }, []);

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
                <MenuList>
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
