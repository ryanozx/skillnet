import React from "react";
import axios from "axios";
import {Button, useToast} from "@chakra-ui/react";
import {PostView} from "./Post";
import { CheckIcon } from '@chakra-ui/icons';

interface EditPostSubmitProps {
    text: string;
    handleClose: () => void;
    postURL: string;
    updatePostHandler: React.Dispatch<React.SetStateAction<PostView>>;
}

export default function EditPostSubmit(props : EditPostSubmitProps) {
    const toast = useToast();

    const onSubmit = () => {
        axios.patch(props.postURL, {"content": props.text}, {withCredentials: true})
        .then(res => {
            props.updatePostHandler(res.data["data"]);
            toast({
                title: "Post updated.",
                description: "Your post has been updated.",
                status: "success",
                duration: 5000,
                isClosable: true,
            });
        })
        .catch(err => {
            console.log(err);
            toast({
                title: "Failed to update post.",
                description: err.response.data.error,
                status: "error",
                duration: 5000,
                isClosable: true,
            });
        })
        .finally(() => {
            props.handleClose();
        })
    }

    return (
    <Button onClick={onSubmit} colorScheme="green" leftIcon={<CheckIcon />} mr={3}>
        Update Post
    </Button>)
}