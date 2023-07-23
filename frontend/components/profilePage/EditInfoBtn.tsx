import React, { useState } from 'react';
import {
    IconButton
} from '@chakra-ui/react';
import { EditIcon} from '@chakra-ui/icons';
import EditInfoModal from './EditInfoModal';
import { User, EditableUserInfo } from '../../types';

interface EditInfoBtnProps {
    user: EditableUserInfo;
    setUser: React.Dispatch<React.SetStateAction<User>>;
}

export default function EditInfoBtn ({user, setUser}: EditInfoBtnProps) {
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
                user={user}
                handleOpen={handleOpen}
                handleClose={handleClose}
                isOpen={isOpen}
                setIsOpen={setIsOpen}
                setUser={setUser}/>
        </>
    );

}