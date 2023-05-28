import React, { useState } from "react";
import {
    Box,
    Button,
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalFooter,
    ModalBody,
    ModalCloseButton,
    FormControl,
    FormLabel,
    Input,
    useToast
} from "@chakra-ui/react";
import { AddIcon } from "@chakra-ui/icons";
import axios from "axios";

export default function CreateCommunityBtn() {
    const [isOpen, setIsOpen] = useState(false);
    const [communityName, setCommunityName] = useState("");
    const toast = useToast();

    const handleOpen = () => setIsOpen(true);
    const handleClose = () => setIsOpen(false);

    const handleInputChange = (event) => {
        setCommunityName(event.target.value);
    };

    const handleSave = () => {
        console.log('API call to create new community');
        axios
            .post('your-endpoint', { name: communityName })
            .then((response) => {
                toast({
                title: "Community created.",
                description: "We've created your community for you.",
                status: "success",
                duration: 5000,
                isClosable: true,
                });
                setCommunityName("");
                handleClose();
            })
            .catch((error) => {
                console.error(error);
                setCommunityName("");
                handleClose();
                toast({
                    title: "Community not created.",
                    description: "we encountered an error: " + error.message + ".",
                    status: "error",
                    duration: 5000,
                    isClosable: true,
                });
            });
      };
      

    return (
        <Box>
            <Button shadow='md' w='100%' leftIcon={<AddIcon />} colorScheme="green" onClick={handleOpen}>
                Create Community
            </Button>

            <Modal isOpen={isOpen} onClose={handleClose}>
                <ModalOverlay />
                <ModalContent>
                    <ModalHeader>Create a new Community</ModalHeader>
                    <ModalCloseButton />
                    <ModalBody>
                        <FormControl>
                            <FormLabel>Community Name</FormLabel>
                            <Input placeholder="Enter community name" value={communityName} onChange={handleInputChange} />
                        </FormControl>
                    </ModalBody>

                    <ModalFooter>
                        <Button colorScheme="blue" mr={3} onClick={handleSave}>
                            Save
                        </Button>
                        <Button variant="ghost" onClick={handleClose}>Cancel</Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
        </Box>
    );
}
