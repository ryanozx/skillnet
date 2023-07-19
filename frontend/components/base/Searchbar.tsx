import React, { useState, useRef, useEffect, ChangeEvent, MouseEvent, MouseEventHandler } from 'react';
import {
  Box,
  Input,
  List,
  ListItem,
  Text,
  Link,
  HStack
} from '@chakra-ui/react';
import axios from 'axios';

interface DropdownItem {
    name: string;
    result_type: string;
    score: number;
    url: string;
  }

const DEBOUNCE_DELAY = 300; // 300 milliseconds


function Searchbar() {
    const [searchTerm, setSearchTerm] = useState<string>('');
    const [dropdownItems, setDropdownItems] = useState<DropdownItem[]>([]);
    const [showDropdown, setShowDropdown] = useState<boolean>(false);
    const dropdownRef = useRef<HTMLInputElement>(null);

    let cancelRequest: (() => void) | null = null;
    let debounceTimeout: NodeJS.Timeout | null = null;

    useEffect(() => {
        let timeoutId: NodeJS.Timeout;
        const handleOutsideClick: EventListener = (event: Event) => {
            if (
                dropdownRef.current &&
                !dropdownRef.current.contains(event.target as Node)
            ) {
                timeoutId = setTimeout(() => {
                    setShowDropdown(false);
                }, 200); // 200ms delay, adjust to your need
            }
        };
        document.addEventListener('mousedown', handleOutsideClick);
        return () => {
            clearTimeout(timeoutId);
            document.removeEventListener('mousedown', handleOutsideClick);
        };
    }, []);

    // Create a new instance of Axios for debounce
    const axiosInstance = axios.create();

    // Cancel the previous request before making a new request
    axiosInstance.interceptors.request.use((request) => {
        const { CancelToken } = axios;
        request.cancelToken = new CancelToken((cancel) => (cancelRequest = cancel));
        return request;
    });

    // Cancel the request if the route changes (react-router)
    useEffect(() => () => {
        if (cancelRequest) cancelRequest();
    }, []);

    // Implement debounce for API call in input change handler
    const handleInputChange = (event: ChangeEvent<HTMLInputElement>) => {
        const { value } = event.target;
        setSearchTerm(value);
        if (debounceTimeout) {
            clearTimeout(debounceTimeout);
        }

        debounceTimeout = setTimeout(() => {
            axiosInstance
            .get(process.env.BACKEND_BASE_URL + '/auth/search', { params: { q: value }, withCredentials: true })
            .then((response) => {
                console.log(response.data.data)
                setDropdownItems(response.data.data);
                setShowDropdown(value !== '' && dropdownRef.current === document.activeElement);
            })
            .catch((error) => console.error('Error fetching search results:', error));
        }, DEBOUNCE_DELAY);
    };


    return (
        <Box position="relative" w="100%">
            <Input
                type="text"
                bg="gray.100"
                placeholder="Search"
                value={searchTerm}
                onChange={handleInputChange}
                ref={dropdownRef}
                
            />
            {showDropdown && (
            
                <Box
                    position="absolute"
                    left={0}
                    bg="white"
                    boxShadow="md"
                    rounded="md"
                    zIndex={10}
                    w="100%"
                >
                    <List 
                        spacing={1}
                        w={dropdownRef.current?.offsetWidth}
                    >
                        {Array.isArray(dropdownItems) && dropdownItems.map((item, index) => (
                            <ListItem key={index} px={4} py={2} _hover={{ bg: 'gray.100', cursor: 'pointer' }}>
                                <Link href={item.url}>
                                    <HStack>
                                        <Text>{item.name}</Text>
                                        <Text color="gray.500">{item.result_type}</Text>
                                    </HStack>
                                </Link>
                            </ListItem>
                        ))}

                    </List>
                </Box>
            
            )}
        </Box>

    );
}

export default Searchbar;
