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
import EditInfoModal from './EditInfoModal';

export default function EditInfoBtn () {
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
            <EditInfoModal
                handleOpen={handleOpen}
                handleClose={handleClose}
                isOpen={isOpen}
                setIsOpen={setIsOpen}
            ></EditInfoModal>
        </>
        

        
    );

}