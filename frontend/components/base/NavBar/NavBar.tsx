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

export default function NavBar(props: any) {
    const { profilePic, username } = props;
    const isDesktop = useBreakpointValue({ base: false, lg: true });
    const toast = useToast();

    const [eventSource, setEventSource] = useState<EventSource | null>(null);

    useEffect(() => {
        console.log("this is being ran");
        if (!eventSource) {
            console.log("this is being ran too");
            const source = new EventSource(process.env.BACKEND_BASE_URL + '/auth/notifications', {
                withCredentials: true,
            });
    
            source.onmessage = function (event) {
                // const notification = JSON.parse(event.data);
                console.log("event: " + event.data);
                // Display the notification. For simplicity, let's use toast here.
                toast({
                    title: "New notification",
                    description: event.data,
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
                    <DesktopNav profilePic={profilePic} username={username}/>
                ) : (
                    <MobileNav profilePic={profilePic} username={username}/>
                )}
            </Flex>
        </Box>
    );
}