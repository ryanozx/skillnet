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


// import React, { useRef, ChangeEvent } from 'react';
// import { Avatar, IconButton, Input } from '@chakra-ui/react';

// interface EditPicBtnProps {
//     currentProfilePic: string;
// }

// export default function EditPicBtn({ currentProfilePic }: EditPicBtnProps) {
//     const fileInputRef = useRef<HTMLInputElement>(null);

    // const handleFileInputClick = () => {
    //     if (fileInputRef.current) {
    //         fileInputRef.current.click();
    //     }
    // };

//     const handleFileChange = (event: ChangeEvent<HTMLInputElement>) => {
//         const file = event.target.files && event.target.files[0];
//         // handle the selected file, e.g. upload it
//     };

//     return (
//         <>
//             <IconButton
//                 aria-label="Change profile picture"
//                 icon={<Avatar size="2xl" src={currentProfilePic} />}
//                 onClick={handleFileInputClick}
//                 isRound={true}
//             />

//             <Input
//                 type="file"
//                 accept="image/*"
//                 ref={fileInputRef}
//                 onChange={handleFileChange}
//                 style={{ display: 'none' }}
//             />
//         </>
//     );
// };

// import React, { useState, useRef, ChangeEvent } from 'react';
// import { Avatar, IconButton, Input, Modal, Button, ModalOverlay, ModalContent, ModalFooter, ModalBody, ModalCloseButton, useDisclosure } from '@chakra-ui/react';
// import Cropper from 'react-cropper';
// import 'cropperjs/dist/cropper.css'; // import the CSS for Cropper

// interface CropperElement extends HTMLImageElement {
//     cropper: {
//         getCroppedCanvas: () => HTMLCanvasElement;
//     };
// }
  
//     const CropperComponent: React.FC = () => {
//     const { isOpen, onOpen, onClose } = useDisclosure();
//     const fileInputRef = useRef<HTMLInputElement>(null);
//     const cropperRef = useRef<CropperElement>(null);
//     const [currentProfilePic, setCurrentProfilePic] = useState<string>("");
//     const [selectedImage, setSelectedImage] = useState<string | undefined>(undefined);
  
//     const handleFileInputClick = () => {
//         if (fileInputRef.current) {
//             fileInputRef.current.click();
//         }
//     };
  
//     const handleFileChange = (event: ChangeEvent<HTMLInputElement>) => {
//         const file = event.target.files && event.target.files[0];
//         if (file && (file.type === 'image/jpeg' || file.type === 'image/png')) {
//             setSelectedImage(URL.createObjectURL(file));
//             onOpen();
//         } else if (file){
//             alert('Please upload a jpg or png file.');
//         }
//     };
  
//     const saveCroppedImage = () => {
//         const imageElement = cropperRef?.current;
//         if (imageElement?.cropper) {
//             setCurrentProfilePic(imageElement.cropper.getCroppedCanvas().toDataURL());
//         }
//         onClose();
//     };

//     return (
//         <>
//             <IconButton
//                 aria-label="Change profile picture"
//                 icon={<Avatar size="2xl" src={currentProfilePic} />}
//                 onClick={handleFileInputClick}
//                 isRound={true}
//             />

//             <Input
//                 type="file"
//                 accept="image/*"
//                 ref={fileInputRef}
//                 onChange={handleFileChange}
//                 style={{ display: 'none' }}
//             />

//             <Modal isOpen={isOpen} onClose={onClose}>
//                 <ModalOverlay />
//                 <ModalContent>
//                 <ModalCloseButton />
//                 <ModalBody>
//                     <Cropper
//                         ref={cropperRef}
//                         src={selectedImage}
//                         style={{ height: 400, width: '100%' }}
//                         aspectRatio={1}
//                         guides={false}
//                     />
//                 </ModalBody>

//                 <ModalFooter>
//                     <Button colorScheme="blue" mr={3} onClick={saveCroppedImage}>
//                         Save
//                     </Button>
//                     <Button variant="ghost" onClick={onClose}>Cancel</Button>
//                 </ModalFooter>
//                 </ModalContent>
//             </Modal>
//         </>
//     );
// }

// export default CropperComponent;



