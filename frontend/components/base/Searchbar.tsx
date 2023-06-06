import React, { useState, useRef, useEffect, ChangeEvent, MouseEvent, MouseEventHandler } from 'react';
import {
  Box,
  Input,
  List,
  ListItem,
  Divider,
  Text,
  VStack,
  Portal
} from '@chakra-ui/react';

function Searchbar() {
    const [searchTerm, setSearchTerm] = useState<string>('');
    const [dropdownItems, setDropdownItems] = useState<string[]>([]);
    // Take note, at some point in time, this might change to useState<ReactNode[]>
    const [showDropdown, setShowDropdown] = useState<boolean>(false);
    const dropdownRef = useRef<HTMLInputElement>(null);

    useEffect(() => {
        const handleOutsideClick: EventListener = (event: Event) => {
            if (
                dropdownRef.current &&
                !dropdownRef.current.contains(event.target as Node)
            ) {
                setShowDropdown(false);
            }
        };
        document.addEventListener('mousedown', handleOutsideClick);
        return () => {
            document.removeEventListener('mousedown', handleOutsideClick);
        };
    }, []);


  const handleInputChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { value } = event.target;
    setSearchTerm(value);

    // Simulating API call to fetch dropdown items based on search term
    // Replace this with your own API logic to fetch data
    const fetchedDropdownItems = [
        'Item 1',
        'Item 2',
        'Item 3',
        'Item 4',
        'Item 5',
    ];

    setDropdownItems(fetchedDropdownItems);
    setShowDropdown(value !== "" && dropdownRef.current === document.activeElement);
  };

  const handleItemClick = (item: string) => {
    // Perform action when dropdown item is clicked
    console.log(`Clicked on ${item}`);
    setShowDropdown(false);
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
                    {dropdownItems.map((item) => (
                        <ListItem
                            key={item}
                            px={4}
                            py={2}
                            _hover={{ bg: 'gray.100', cursor: 'pointer' }}
                            onClick={() => handleItemClick(item)}
                        >
                        <Text>{item}</Text>
                        </ListItem>
                    ))}
                </List>
            </Box>
        
        )}
    </Box>

  );
}

export default Searchbar;
