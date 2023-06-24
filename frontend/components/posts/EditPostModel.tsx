import React, {useEffect, useState, useRef} from "react"
import {PostComponent, PostView} from "./Post"
import {Modal, ModalContent, ModalFooter, ModalHeader, ModalOverlay, Textarea, useToast} from "@chakra-ui/react";
import EditPostSubmit from "./EditPostSubmit";
import EditPostCancel from "./EditPostCancel";

interface EditPostModalProps {
    isOpen: boolean;
    setIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
    postComponent: PostComponent;
    updatePostHandler: React.Dispatch<React.SetStateAction<PostView>>;
}

export default function EditPostModel(props : EditPostModalProps) {
    const [text, setText] = useState<string>("");
    const didMountRef = useRef(false);
    const toast = useToast();
    
    const handleClose = () => props.setIsOpen(false);
    const base_url = process.env.BACKEND_BASE_URL;
    const postURL = base_url + `/auth/posts/${props.postComponent.ID}`;

    const loadText = () => {
        setText(props.postComponent.Content);
    };

    useEffect(() => {
        if (didMountRef.current) {
            loadText();
        }
        didMountRef.current = true;
    }, [props.isOpen]);

    return <Modal isOpen={props.isOpen} onClose={handleClose} size={{ base: 'md', md: '2xl' }} closeOnOverlayClick={false}>
        <ModalOverlay />
        <ModalContent padding="20px">
            <ModalHeader paddingInline="0px">Edit Post</ModalHeader>
            <Textarea size="md" value={text} onChange={e => setText(e.target.value)}/>
            <ModalFooter display="space-between" paddingInline="0px">
                <EditPostSubmit text={text} handleClose={handleClose} postURL={postURL} updatePostHandler={props.updatePostHandler}/>
                <EditPostCancel handleClose={handleClose} />
            </ModalFooter>
        </ModalContent>
    </Modal>
}