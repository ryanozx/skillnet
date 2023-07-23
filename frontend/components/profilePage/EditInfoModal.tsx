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
import { escapeHtml } from '../../types';

export interface FormType {
    name: string;
    title: string;
    about: string;
    privacySettings: {
      [key: string]: boolean;
    };
}

interface EditProfileModalProps {
    user: User;
    setUser: React.Dispatch<React.SetStateAction<User>>;
    handleOpen: () => void;
    handleClose: () => void;
    isOpen: boolean;
    setIsOpen: (isOpen: boolean) => void;
}  

export default function EditProfileModal(props: EditProfileModalProps) {
    const { handleClose, isOpen, setIsOpen, setUser } = props;
    const toast = useToast();
    const [activeTab, setActiveTab] = useState(0);
    const [form, setForm] = useState<FormType>({
        name: props.user.Name,
        title: props.user.Title,
        about: props.user.AboutMe,
        privacySettings: {
            title: props.user.ShowTitle,
            about: props.user.ShowAboutMe,
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

    const handleSave = () => {
        const base_url = process.env.BACKEND_BASE_URL;
        const url = base_url + "/auth/user"
        axios.patch(url, {
                "Name": escapeHtml(form.name),
                "Title": escapeHtml(form.title),
                "AboutMe": escapeHtml(form.about),
                "ShowAboutMe": form.privacySettings["about"],
                "ShowTitle": form.privacySettings["title"],
            }, {
                withCredentials: true,
            })
            .then(res => {
                setUser({...res.data.data});
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
        <>
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
        </>
    );
};

