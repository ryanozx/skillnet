import React, {useRef} from "react"
import axios from "axios"
import {AlertDialog, AlertDialogBody, AlertDialogContent, AlertDialogFooter, AlertDialogHeader, AlertDialogOverlay, Button, MenuItem, useToast, useDisclosure} from "@chakra-ui/react"
import { PostComponent } from "./Post";

interface DeletePostItemProps {
    post: PostComponent;
    deleteHandler: React.Dispatch<React.SetStateAction<boolean>>;
}

export default function DeletePostItem(props : DeletePostItemProps) {
    const {isOpen, onOpen, onClose} = useDisclosure();
    const cancelRef = useRef<HTMLButtonElement>(null);
    const postURL = `http://localhost:8080/auth/posts/${props.post.ID}`;
    const toast = useToast();

    const handleDelete = () => {
        axios.delete(postURL, {withCredentials: true})
        .then(res => {
            props.deleteHandler(true);
            toast({
                title: "Post deleted.",
                description: "Your post has been deleted.",
                status: "success",
                duration: 5000,
                isClosable: true,
            });
        })
        .catch(err => {
            console.log(err);
            toast({
                title: "An error occurred.",
                description: err.response.data.error,
                status: "error",
                duration: 5000,
                isClosable: true,
            });
        })
        .finally(() => {
            onClose();
        })
    }

    return (
        <>
            <MenuItem
                onClick={onOpen}
            >Delete post
            </MenuItem>

            <AlertDialog
                isOpen={isOpen}
                leastDestructiveRef={cancelRef}
                onClose={onClose}
            >
                <AlertDialogOverlay>
                    <AlertDialogContent>
                        <AlertDialogHeader fontSize="lg" fontWeight="bold">
                            Delete Post
                        </AlertDialogHeader>
                        
                        <AlertDialogBody>
                            Are you sure? You cannot undo this action afterwards.
                        </AlertDialogBody>

                        <AlertDialogFooter>
                            <Button ref={cancelRef} onClick={onClose}>
                                Cancel
                            </Button>
                            <Button colorScheme="red" onClick={handleDelete} ml={3}>
                                Delete
                            </Button>
                        </AlertDialogFooter>
                    </AlertDialogContent>
                </AlertDialogOverlay>
            </AlertDialog>
        </>
    )
}