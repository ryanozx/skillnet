import React, { useState, useEffect } from 'react';
import axios from 'axios';  
import { 
    Button,  
    Modal, 
    ModalContent, 
    ModalFooter, 
    ModalOverlay, 
    Tab,
    Tabs,
    TabList,
    TabPanel,
    TabPanels,
    useToast
} from "@chakra-ui/react";
import { CloseIcon, CheckIcon } from '@chakra-ui/icons';
import BasicInfoForm from './BasicInfoForm';
import PrivacyForm from './PrivacyForm';
import { User } from '../../types';

interface FormType {
    name: string;
    title: string;
    about: string;
    privacySettings: {
      [key: string]: boolean;
    };
}

interface EditProfileModalProps {
    setUser: React.Dispatch<React.SetStateAction<User>>;
    handleOpen: () => void;
    handleClose: () => void;
    isOpen: boolean;
    setIsOpen: (isOpen: boolean) => void;
}  

export default function EditProfileModal(props: EditProfileModalProps) {
    const { handleOpen, handleClose, isOpen, setIsOpen, setUser } = props;
    const toast = useToast();
    const [activeTab, setActiveTab] = useState(0);
    const [form, setForm] = useState<FormType>({
        name: "",
        title: "",
        about: "",
        privacySettings: {
            tagline: false,
            about: false,
            projects: false,
            activity: false
        }
    });

    const handleInputChange = (e: any) => {
        const { name, value } = e.target;
        setForm(prevState => ({
            ...prevState,
            [name]: value
        }));
    };

    const handleSwitchChange = (name: string) => {
        setForm(prevState => ({
            ...prevState,
            privacySettings: {
                ...prevState.privacySettings,
                [name]: !prevState.privacySettings[name]
            }
        }));
    };

    const handleTabChange = (index: any) => {
        setActiveTab(index);
    };

    useEffect(() => {
        if (isOpen) {
            const sessionId = sessionStorage.getItem('sessionId');
            console.log('API call to get the current privacy setting of user');
            axios
                .get('your-privacy-endpoint', {
                headers: {
                    Authorization: `Bearer ${sessionId}`
                }
                })
                .then(response => {
                setForm(prevState => ({
                    ...prevState,
                    privacySettings: response.data
                }));
                })
                .catch(error => {
                console.error(error);
                });
        }
        
    }, [isOpen]);

    const handleSave = () => {
        const url = "http://localhost:8080/auth/user"
        axios.patch(url, {
                name: form.name,
                title: form.title,
                aboutme: form.about,
            }, {
                withCredentials: true,
            })
            .then(res => {
                console.log(res.data); // Log the response data
                const { AboutMe, Name, Title } = res.data.data;
                setUser((prevUser: User) => ({
                    ...prevUser,
                    AboutMe: AboutMe ? AboutMe : "No description available",
                    Name: Name ? Name :  "No display name",
                    Title: Title ? Title : "No title available",
                }));
                toast({
                    title: "Profile updated.",
                    description: "We've updated your profile for you.",
                    status: "success",
                    duration: 9000,
                    isClosable: true,
                });
                setIsOpen(false);
            })
            .catch(error => {
                console.error(error);
            });
    };

    return (
        
        <Modal isOpen={isOpen} onClose={handleClose} size={{ base: 'md', md: '2xl' }} closeOnOverlayClick={false}>
            <ModalOverlay />
            <ModalContent>
                <Tabs index={activeTab} onChange={handleTabChange}>
                <TabList>
                    <Tab>Basic</Tab>
                    <Tab>Privacy</Tab>
                </TabList>
                <TabPanels>
                    <TabPanel>
                        <BasicInfoForm form={form} handleInputChange={handleInputChange} />
                    </TabPanel>
                    <TabPanel>
                        <PrivacyForm form={form} handleSwitchChange={handleSwitchChange} />
                    </TabPanel>
                </TabPanels>
                </Tabs>
                <ModalFooter>
                    <Button onClick={handleSave} colorScheme="green" leftIcon={<CheckIcon />} mr={3}>
                        Save
                    </Button>
                    <Button colorScheme="red" onClick={handleClose} leftIcon={<CloseIcon />}>
                        Cancel
                    </Button>                
                </ModalFooter>
            </ModalContent>
            </Modal>

    );
};

