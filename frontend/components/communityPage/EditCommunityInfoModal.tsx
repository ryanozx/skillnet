import React, {useEffect, useState, useRef} from "react"
import {Button, Modal, ModalContent, ModalFooter, ModalHeader, ModalOverlay, Textarea, useToast} from "@chakra-ui/react";

import {Community} from "./CommunityInfo";
import {CloseIcon, CheckIcon} from "@chakra-ui/icons";
import axios from "axios";
import { escapeHtml } from "../../types";

interface EditCommunityModalProps {
    isOpen: boolean;
    setIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
    community: Community;
    updateCommunityHandler: React.Dispatch<React.SetStateAction<Community>>;
}

export default function EditCommunityModal(props : EditCommunityModalProps) {
    const [text, setText] = useState<string>("");
    const didMountRef = useRef(false);
    
    const handleClose = () => props.setIsOpen(false);
    const base_url = process.env.BACKEND_BASE_URL;
    const communityURL = base_url + `/auth/community/${props.community.Name}`;
    const toast = useToast();

    const loadText = () => {
        setText(props.community.About);
    };

    const onSubmit = () => {
        axios.patch(communityURL, {"about": escapeHtml(text)}, {withCredentials: true})
        .then(res => {
            props.updateCommunityHandler(res.data["data"]["Community"]);
            toast({
                title: "Community updated.",
                description: "Your community has been updated.",
                status: "success",
                duration: 5000,
                isClosable: true,
            });
        })
        .catch(err => {
            console.log(err);
            toast({
                title: "Failed to update community.",
                description: err.response.data.error,
                status: "error",
                duration: 5000,
                isClosable: true,
            });
        })
        .finally(() => {
            handleClose();
        })
    }

    useEffect(() => {
        if (didMountRef.current) {
            loadText();
        }
        didMountRef.current = true;
    }, [props.isOpen]);

    return <Modal isOpen={props.isOpen} onClose={handleClose} size={{ base: 'md', md: '2xl' }} closeOnOverlayClick={false}>
        <ModalOverlay />
        <ModalContent padding="20px">
            <ModalHeader paddingInline="0px">Edit About Community</ModalHeader>
            <Textarea size="md" value={text} onChange={e => setText(e.target.value)}/>
            <ModalFooter display="space-between" paddingInline="0px">
                <Button onClick={onSubmit} colorScheme="green" leftIcon={<CheckIcon />} mr={3}>
                    Save Changes
                </Button>
                <Button onClick={handleClose} colorScheme="red" leftIcon={<CloseIcon />}>
                    Discard Changes
                </Button>
            </ModalFooter>
        </ModalContent>
    </Modal>
}