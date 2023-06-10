import React, { useState } from 'react';
import { useDisclosure } from '@chakra-ui/react';
import EditPicButton from './EditPicButton';
import ImageCropper from './ImageCropper';

interface CropperComponentProps {
    profilePic: string | undefined;
}

const CropperComponent: React.FC<CropperComponentProps> = ({profilePic}) => {
    const { isOpen, onOpen, onClose } = useDisclosure();
    const [currentProfilePic, setCurrentProfilePic] = useState<string>(profilePic || "");
    const [selectedImage, setSelectedImage] = useState<string | undefined>(undefined);
  
    const handleValidFile = (file: File) => {
        setSelectedImage(URL.createObjectURL(file));
        onOpen();
    };
  
    const handleCroppedImage = (dataUrl: string) => {
        setCurrentProfilePic(dataUrl);
    };

    return (
        <>
            <EditPicButton currentProfilePic={currentProfilePic} onValidFile={handleValidFile} />
            <ImageCropper isOpen={isOpen} onClose={onClose} onCropped={handleCroppedImage} imageSrc={selectedImage} />
        </>
    );
};

export default CropperComponent;
