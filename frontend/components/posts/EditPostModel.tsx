import React, {useEffect, useState} from "react"
import {PostComponent} from "./Post"
import {Modal, ModalContent, ModalFooter, ModalHeader, ModalOverlay, Textarea, useToast} from "@chakra-ui/react";
import axios from "axios";

interface EditPostModalProps {
    isOpen: boolean;
    setIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
    postComponent: PostComponent;
}

export default function EditPostModel(props : EditPostModalProps) {
    const [text, setText] = useState<string>("");
    const toast = useToast();
    
    const handleClose = () => props.setIsOpen(false);
    const postURL = `http://localhost:8080/auth/posts/${props.postComponent.ID}`;

    const loadText = async() => {
        console.log(postURL);
        axios.get(postURL, {withCredentials: true})
        .then((response) => {
            console.log(response.data["data"]["Post"]["Content"]);
            setText(response.data["data"]["Post"]["Content"]);
        })
        .catch((error) => {
            console.log(error);
            toast({
                title: "An error occured",
                description: error.message,
                status: "error",
                duration: 5000,
                isClosable: true,
            });
        })
    };

    useEffect(() => {
        loadText();
    }, []);

    return <Modal isOpen={props.isOpen} onClose={handleClose} size={{ base: 'md', md: '2xl' }}>
        <ModalOverlay />
        <ModalContent padding="20px">
            <ModalHeader paddingInline="0px">Edit Post</ModalHeader>
            <Textarea size="md">

            </Textarea>
            <ModalFooter>

            </ModalFooter>
        </ModalContent>
    </Modal>
}