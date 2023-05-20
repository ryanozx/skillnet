import React, { useState } from 'react';
import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalCloseButton,
  ModalBody,
  ModalFooter,
  FormControl,
  FormLabel,
  Input,
  Textarea,
  IconButton,
  Button,
} from '@chakra-ui/react';
import { EditIcon, CheckIcon, CloseIcon } from '@chakra-ui/icons';

export default function EditProfileModal() {
    const [isOpen, setIsOpen] = useState(false);

    const handleOpen = () => setIsOpen(true);
    const handleClose = () => setIsOpen(false);

    return (
        <>
        <IconButton 
            alignSelf={"flex-end"}
            icon={<EditIcon />} 
            onClick={handleOpen} 
            aria-label="Edit profile"
        />

        <Modal isOpen={isOpen} onClose={handleClose} size="xl">
            <ModalOverlay />
            <ModalContent>
                <ModalHeader>Edit your profile</ModalHeader>
                <ModalCloseButton />
                <ModalBody>
                    <FormControl id="name">
                        <FormLabel>Name</FormLabel>
                        <Input placeholder="Your name" />
                    </FormControl>
                    <FormControl id="title" mt={4}>
                        <FormLabel>Title</FormLabel>
                        <Input placeholder="Your title" />
                    </FormControl>
                    <FormControl id="about" mt={4}>
                    <FormLabel>About me</FormLabel>
                    <Textarea placeholder="About you" />
                    </FormControl>
                </ModalBody>

                <ModalFooter>
                    <Button mr={3} colorScheme="red" onClick={handleClose} leftIcon={<CloseIcon />}>
                        Cancel
                    </Button>
                    <Button colorScheme="green" leftIcon={<CheckIcon />}>
                        Save
                    </Button>
                </ModalFooter>
            </ModalContent>
        </Modal>
        </>
    );
};
