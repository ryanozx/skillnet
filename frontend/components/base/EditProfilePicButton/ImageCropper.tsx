import React, { useRef } from 'react';
import { Button, Modal, ModalOverlay, ModalContent, ModalBody, ModalFooter, ModalCloseButton, useDisclosure } from '@chakra-ui/react';
import axios from 'axios';
import Cropper from 'react-cropper';
import 'cropperjs/dist/cropper.css';

interface CropperElement extends HTMLImageElement {
    cropper: {
        getCroppedCanvas: () => HTMLCanvasElement;
    };
}

interface ImageCropperProps {
    isOpen: boolean;
    onClose: () => void;
    onCropped: (imageDataUrl: string) => void;
    imageSrc: string | undefined;
}

const ImageCropper: React.FC<ImageCropperProps> = ({ isOpen, onClose, onCropped, imageSrc }) => {
    const cropperRef = useRef<CropperElement>(null);

    const saveCroppedImage = async () => {
        const imageElement = cropperRef?.current;
        const url = "http://localhost:8080/auth/user/photo";
        if (imageElement?.cropper) {
            const croppedCanvas = imageElement.cropper.getCroppedCanvas();
            croppedCanvas.toBlob(async (blob) => {
                if (blob) {
                    const formData = new FormData();
                    formData.append('file', blob, 'image.jpg');
                    axios.post(url, formData, {
                        headers: {
                            'Content-Type': 'multipart/form-data',
                        }, withCredentials: true,
                    })
                    .then((response) => {
                        console.log(response)
                        if (response.data.url) {
                            onCropped(response.data.url);
                        }
                    })
                    .catch((error) => {
                        console.error('Failed to upload the cropped image', error);
                    })
                }
            }, 'image/jpeg');
            onClose();
        }
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose}>
            <ModalOverlay />
            <ModalContent>
                <ModalCloseButton />
                <ModalBody>
                    <Cropper
                        ref={cropperRef}
                        src={imageSrc}
                        style={{ height: 400, width: '100%' }}
                        aspectRatio={1}
                        guides={false}
                    />
                </ModalBody>

                <ModalFooter>
                    <Button colorScheme="blue" mr={3} onClick={saveCroppedImage}>
                        Save
                    </Button>
                    <Button variant="ghost" onClick={onClose}>Cancel</Button>
                </ModalFooter>
            </ModalContent>
        </Modal>
    );
};

export default ImageCropper;
