import { useState } from 'react';
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

export default function EditProfileModal(props) {
    const { handleOpen, handleClose, isOpen, setIsOpen } = props;
    
    const [activeTab, setActiveTab] = useState(0);
    const [form, setForm] = useState({
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

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        setForm(prevState => ({
            ...prevState,
            [name]: value
        }));
    };

    const handleSwitchChange = (name) => {
        setForm(prevState => ({
            ...prevState,
            privacySettings: {
                ...prevState.privacySettings,
                [name]: !prevState.privacySettings[name]
            }
        }));
    };


    const handleTabChange = (index) => {
        setActiveTab(index);
    };

    const handleSave = () => {
        console.log(form); // Log the form data. Here is where you could add your API request to update the user data.
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

