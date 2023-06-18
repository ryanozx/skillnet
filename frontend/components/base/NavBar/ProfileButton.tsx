import React, { MouseEventHandler } from 'react';
import {
    Avatar,
    Button,
    Link,
    Menu,
    MenuButton,
    MenuList,
    MenuItem,
    MenuDivider,
    useToast
} from '@chakra-ui/react';
import axios from 'axios';
import { useRouter } from 'next/router'; // Import the useRouter hook

interface ProfileButtonProps {
    profilePic: string;
    username: string;
}

export default function ProfileButton(props: ProfileButtonProps) {
    const { profilePic, username } = props;
    const toast = useToast();
    const router = useRouter(); // Use the useRouter hook

    const handleClick: MouseEventHandler = () => {
        var url = 'http://localhost:8080/auth/logout'
        axios.post(url, {}, {withCredentials: true}) // Make a POST request to the logout endpoint
            .then((res) => {
                console.log(res);
                toast({
                    title: "Form submission successful.",
                    description: "We've successfully logged you out.",
                    status: "success",
                    duration: 5000,
                    isClosable: true,
                });
                router.push("/"); // Navigate to "/"
            }).catch((error) => {
                console.log(error.response);
                toast({
                    title: "An error occurred.",
                    description: error.response.data.error,
                    status: "error",
                    duration: 5000,
                    isClosable: true,
                });
            });
    };

    return (
        <Menu>
                            
            <MenuButton
                as={Button}
                rounded={'full'}
                variant={'link'}
                cursor={'pointer'}
                minW={0}
            >
                <Avatar
                size={'md'}
                src={profilePic}
                />
            </MenuButton>
            <MenuList>
                <MenuItem><Link href={`/profile/${username}`}>View your profile</Link></MenuItem>
                <MenuDivider />
                <MenuItem onClick={handleClick}>logout</MenuItem> {/* OnClick of logout, call handleClick */}
            </MenuList>
        </Menu>
    );
}
