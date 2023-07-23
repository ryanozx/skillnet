import React, {useEffect, useState, useRef} from "react"
import {Button, Modal, ModalContent, ModalFooter, ModalHeader, ModalOverlay, Textarea, useToast} from "@chakra-ui/react";

import {CloseIcon, CheckIcon} from "@chakra-ui/icons";
import axios from "axios";
import { Project } from "./ProjectInfo";
import { escapeHtml } from "../../types";

interface EditProjectInfoModalProps {
    isOpen: boolean,
    setIsOpen: React.Dispatch<React.SetStateAction<boolean>>,
    project: Project,
    updateProjectHandler: React.Dispatch<React.SetStateAction<Project>>,
}

export default function EditProjectInfoModal(props : EditProjectInfoModalProps) {
    const [text, setText] = useState<string>("");
    const didMountRef = useRef(false);
    
    const handleClose = () => props.setIsOpen(false);
    const base_url = process.env.BACKEND_BASE_URL;
    const projectURL = base_url + `/auth/projects/${props.project.ID}`;
    const toast = useToast();

    const loadText = () => {
        setText(props.project.ProjectInfo);
    };

    const onSubmit = () => {
        axios.patch(projectURL, {"projectInfo": escapeHtml(text)}, {withCredentials: true})
        .then(res => {
            console.log(res.data.data)
            props.updateProjectHandler({...res.data.data});
            toast({
                title: "Project updated.",
                description: "Your project has been updated.",
                status: "success",
                duration: 5000,
                isClosable: true,
            });
        })
        .catch(err => {
            console.log(err);
            toast({
                title: "Failed to update project.",
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
            <ModalHeader paddingInline="0px">Edit About Project</ModalHeader>
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