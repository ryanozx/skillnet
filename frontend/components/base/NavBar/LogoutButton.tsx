import React, { MouseEventHandler } from 'react';
import {
    MenuItem,
    useToast
} from '@chakra-ui/react';
import axios from 'axios';
import { useRouter } from 'next/router'; // Import the useRouter hook

const LogoutButton: React.FC = () => {
    const toast = useToast();
    const router = useRouter();

    const handleClick: MouseEventHandler = () => {
        const base_url = process.env.BACKEND_BASE_URL;
        axios.post(base_url + '/auth/logout', {}, {withCredentials: true})
            .then((res) => {
                console.log(res);
                toast({
                    title: "Form submission successful.",
                    description: "We've successfully logged you out.",
                    status: "success",
                    duration: 5000,
                    isClosable: true,
                });
                router.push("/");
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
        <MenuItem onClick={handleClick}>Log out</MenuItem> 
    )
}
export default LogoutButton;
