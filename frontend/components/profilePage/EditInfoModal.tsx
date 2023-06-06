import React, { useState, useEffect } from 'react';
import axios from 'axios';  
import { 
    Button, 
    FormControl, 
    FormLabel, 
    Modal, 
    ModalBody, 
    ModalCloseButton, 
    ModalContent, 
    ModalFooter, 
    ModalHeader, 
    ModalOverlay, 
    Input, 
    Textarea, 
    Switch,
    Tab,
    Tabs,
    TabList,
    TabPanel,
    TabPanels
} from "@chakra-ui/react";
import { CloseIcon, CheckIcon } from '@chakra-ui/icons';
import BasicInfoForm from './BasicInfoForm';
import PrivacyForm from './PrivacyForm';

interface FormType {
    name: string;
    title: string;
    about: string;
    privacySettings: {
      [key: string]: boolean;
    };
}
  

export default function EditProfileModal(props: any) {
    const { handleOpen, handleClose, isOpen, setIsOpen } = props;
    
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
        const sessionId = sessionStorage.getItem('sessionId');
        console.log('API call to save updated changes to profile page');
        axios
            .post('your-update-endpoint', form, {
                headers: {
                Authorization: `Bearer ${sessionId}`
                }
            })
            .then(response => {
                // console.log(response.data); // Log the response data
                console.log('Successfully updated');
                setIsOpen(false);
            })
            .catch(error => {
                console.error(error);
            });
    };

    return (
        
        <Modal isOpen={isOpen} onClose={handleClose} size={{ base: 'md', md: '2xl' }}>
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

