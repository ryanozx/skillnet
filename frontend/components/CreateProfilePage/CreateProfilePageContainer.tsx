import React, { useState, useEffect } from 'react';
import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  Text,
  Textarea,
  Flex,
  Stack,
  HStack,
  VStack,
  Divider
} from '@chakra-ui/react';
import axios from 'axios';
import CropperComponent from './CropperComponent';
import { useRouter } from 'next/router';

const NameTitleForm: React.FC<{form: any, handleChange: (e: React.ChangeEvent<HTMLInputElement>) => void}> = ({form, handleChange}) => {
    return (
        <>
            <Stack>
                <FormControl>
                    <FormLabel>Name</FormLabel>
                    <Input
                        type="text"
                        name="name"
                        value={form.name}
                        onChange={handleChange}
                    />
                </FormControl>
                <FormControl>
                    <FormLabel>Title</FormLabel>
                    <Input
                        type="text"
                        name="title"
                        value={form.title}
                        onChange={handleChange}
                    />
                </FormControl>
            </Stack>
            
        </>
    );
}

const AboutMeForm: React.FC<{form: any, handleChange: (e: React.ChangeEvent<HTMLTextAreaElement>) => void}> = ({form, handleChange}) => {
    return (
        <FormControl>
            <FormLabel>About Me</FormLabel>
            <Textarea
                name="aboutMe"
                value={form.aboutMe}
                onChange={handleChange}
            />
        </FormControl>
    );
}

const FormButtons: React.FC<{handleSubmit: () => void}> = ({handleSubmit}) => {
    const router = useRouter();
    const handleSkip = () => {
        router.push(`/profile/me`)
    }

    return (
        <Stack isInline justifyContent="flex-end">
            <Button variant="outline" mr={2} onClick={handleSkip}>
                Skip
            </Button>
            <Button colorScheme="teal" onClick={handleSubmit}>
                Save
            </Button>
        </Stack>
    );
}


const CreateProfilePageContainer: React.FC = () => {
    const [form, setForm] = useState({
        aboutMe: '',
        name: '',
        title: '',
        profilePic: '',
    });

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
                console.log(res.data);
            })
            .catch((error) => {
                console.log(error);
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
                            <NameTitleForm form={form} handleChange={handleChange}/>
                        </Box>
                    </HStack>
                    <AboutMeForm form={form} handleChange={handleChange}/>
                    <Divider/>
                    <FormButtons handleSubmit={handleSubmit}/>
                </VStack>
            </Box>
        </Flex>
    </>
    

  );
};

export default CreateProfilePageContainer;
