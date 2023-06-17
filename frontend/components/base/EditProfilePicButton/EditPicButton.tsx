// ImageUploadButton.tsx
import React, { ChangeEvent, useRef } from 'react';
import { Avatar, IconButton, Input } from '@chakra-ui/react';

interface EditPicButtonProps {
    onValidFile: (file: File) => void;
    currentProfilePic: string;
}

const EditPicButton: React.FC<EditPicButtonProps> = ({ onValidFile, currentProfilePic }) => {
    const fileInputRef = useRef<HTMLInputElement>(null);

    const handleFileInputClick = () => {
        fileInputRef.current?.click();
    };

    const handleFileChange = (event: ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files && event.target.files[0];
        if (file && (file.type === 'image/jpeg' || file.type === 'image/png')) {
            onValidFile(file);
        } else if (file) {
            alert('Please upload a jpg or png file.');
        }
    };

    return (
        <>
            <IconButton
                aria-label="Change profile picture"
                icon={<Avatar size="2xl" src={currentProfilePic} />}
                onClick={handleFileInputClick}
                isRound={true}
            />
            <Input
                type="file"
                accept="image/*"
                ref={fileInputRef}
                onChange={handleFileChange}
                style={{ display: 'none' }}
            />
        </>
    );
};

export default EditPicButton;



