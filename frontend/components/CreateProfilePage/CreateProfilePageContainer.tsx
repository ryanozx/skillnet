import React, { useState, useEffect } from 'react';
import {
  Box,
  Text,
  Flex,
  HStack,
  VStack,
  Divider,
  useToast
} from '@chakra-ui/react';
import axios from 'axios';
import CropperComponent from './CropperComponent';
import NameTitleFields from './NameTitleFields';
import AboutMeField from './AboutMeField';
import FormButtons from './FormButtons';
import { requireAuth } from '../../withAuthRedirect';
import {useRouter} from "next/router";

export default requireAuth(function CreateProfilePageContainer() {
    const [form, setForm] = useState({
        aboutMe: '',
        name: '',
        title: '',
        profilePic: '',
    });

    const router = useRouter();
    const toast = useToast();

    useEffect(() => {
        axios
        .get('http://localhost:8080/auth/user', { withCredentials: true })
        .then((res) => {
            const { AboutMe, Name, Title, ProfilePic } = res.data.data;
            setForm({
                aboutMe: AboutMe ? AboutMe : '',
                name: Name ? Name : '',
                title: Title ? Title : '',
                profilePic: ProfilePic ? ProfilePic : '',
            });
        })
        .catch((error) => {
            console.log(error);
        });
    }, []);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        setForm({ ...form, [e.target.name]: e.target.value });
    };

    const handleSubmit = () => {
        axios.patch('http://localhost:8080/auth/user', form, { withCredentials: true })
            .then((res) => {
                toast({
                    title: "Form submission successful.",
                    description: "Your profile has been updated!",
                    status: "success",
                    duration: 5000,
                    isClosable: true,
                });
                router.push("/feed");
            })
            .catch((error) => {
                console.log(error);
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
        <>
            <Flex justify="center" align="center" h="100vh" direction={{base: "column", lg: "row"}}>
                <Box 
                    p={5} shadow="md" borderWidth="1px" borderRadius="md" 
                    w={{base: "90vw", lg: "60vw"}}>
                    <VStack spacing={10} alignItems="stretch">
                        <Text fontSize="3xl" fontWeight="bold">Create Your Profile</Text>
                        <HStack pl={{base: 0, md:20}} justify="center">
                            <Box flex={2}>
                                <CropperComponent profilePic={form.profilePic}/>
                            </Box>
                            <Box flex={3}>
                                <NameTitleFields form={form} handleChange={handleChange}/>
                            </Box>
                        </HStack>
                        <AboutMeField form={form} handleChange={handleChange}/>
                        <Divider/>
                        <FormButtons handleSubmit={handleSubmit}/>
                    </VStack>
                </Box>
            </Flex>
        </>
    );
});

// export default requireAuth(CreateProfilePageContainer);
