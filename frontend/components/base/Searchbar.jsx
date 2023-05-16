import { useState, useRef, useEffect } from 'react';
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
  const [searchTerm, setSearchTerm] = useState('');
  const [dropdownItems, setDropdownItems] = useState([]);
  const [showDropdown, setShowDropdown] = useState(false);
  const dropdownRef = useRef(null);

useEffect(() => {
    const handleOutsideClick = (event) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target)
      ) {
        setShowDropdown(false);
      }
    };

    document.addEventListener('mousedown', handleOutsideClick);

    return () => {
      document.removeEventListener('mousedown', handleOutsideClick);
    };
  }, []);

  const handleInputChange = (event) => {
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
    setShowDropdown(value && dropdownRef.current === document.activeElement);
  };

  const handleItemClick = (item) => {
    // Perform action when dropdown item is clicked
    console.log(`Clicked on ${item}`);
    setShowDropdown(false);
  };

  return (
    <Box position="relative">
      <Input
        type="text"
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
          >
            <List 
                spacing={1}
                w={dropdownRef.current?.offsetWidth}
                minW="100%"
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
