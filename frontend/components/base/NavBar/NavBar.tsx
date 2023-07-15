import {
    Box,
    Flex,
    useBreakpointValue,
    useToast,
} from '@chakra-ui/react';
import DesktopNav from './DesktopNav';
import MobileNav from './MobileNav';
import React, { useState, useEffect } from 'react';

interface NavBarProps {
    profilePic: string;
    username: string;
}

export default function NavBar(props: NavBarProps) {
    const { profilePic, username } = props;
    const isDesktop = useBreakpointValue({ base: false, lg: true });
    const toast = useToast();

    const [eventSource, setEventSource] = useState<EventSource | null>(null);
    const [notifications, setNotifications] = useState<string[]>([]);
    const [hasNewNotifications, setHasNewNotifications] = useState<boolean>(false);

    useEffect(() => {
        if (!eventSource) {
            const source = new EventSource(process.env.BACKEND_BASE_URL + '/auth/notifications', {
                withCredentials: true,
            });
    
            source.onmessage = function (event) {
                const notification = JSON.parse(event.data)

                setNotifications(prevNotifications => [notification.content, ...prevNotifications]);
                setHasNewNotifications(true);
                toast({
                    title: "New notification",
                    description: notification.content,  
                    status: "info",
                    duration: 5000,
                    isClosable: true,
                });
            };
    
            source.onerror = function (event) {
                console.error("EventSource failed:", event);
            };
    
            setEventSource(source);
        }
    
        // When the component unmounts, close the event source
        return () => {
            eventSource?.close();
        };
    }, [eventSource]);
    
    
    return (
        <Box>
            <Flex py={2} px={4} borderBottom={1} align={'center'}>
                {isDesktop ? (
                    <DesktopNav 
                        profilePic={profilePic} 
                        username={username} 
                        notifications={notifications}
                        hasNewNotifications={hasNewNotifications}
                        setHasNewNotifications={setHasNewNotifications}
                        />
                ) : (
                    <MobileNav 
                        profilePic={profilePic} 
                        username={username} 
                        notifications={notifications}
                        hasNewNotifications={hasNewNotifications}
                        setHasNewNotifications={setHasNewNotifications}
                        />
                )}
            </Flex>
        </Box>
    );
}
