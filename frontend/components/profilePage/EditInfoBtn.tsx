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
import { User } from '../../types';

interface EditInfoBtnProps {
    setUser: React.Dispatch<React.SetStateAction<User>>;
}

export default function EditInfoBtn ({setUser}: EditInfoBtnProps) {
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
                setUser={setUser}/>
        </>
    );

}