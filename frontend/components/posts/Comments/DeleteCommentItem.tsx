import React, {useRef} from "react"
import axios from "axios"
import {AlertDialog, AlertDialogBody, AlertDialogContent, AlertDialogFooter, AlertDialogHeader, AlertDialogOverlay, Button, MenuItem, useToast, useDisclosure} from "@chakra-ui/react"
import { CommentComponent } from "./Comment";

interface DeleteCommentItemProps {
    comment: CommentComponent;
    deleteHandler: React.Dispatch<React.SetStateAction<boolean>>;
    commentCountHandler: React.Dispatch<React.SetStateAction<number>>;
}

export default function DeleteCommentItem(props : DeleteCommentItemProps) {
    const {isOpen, onOpen, onClose} = useDisclosure();
    const cancelRef = useRef<HTMLButtonElement>(null);
    const base_url = process.env.BACKEND_BASE_URL;
    const postURL = base_url + `/auth/comments/${props.comment.ID}`;
    const toast = useToast();

    const handleDelete = () => {
        axios.delete(postURL, {withCredentials: true})
        .then(res => {
            props.deleteHandler(true);
            props.commentCountHandler(res.data["data"]["CommentCount"])
            toast({
                title: "Comment deleted.",
                description: "Your comment has been deleted.",
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
            >Delete comment
            </MenuItem>

            <AlertDialog
                isOpen={isOpen}
                leastDestructiveRef={cancelRef}
                onClose={onClose}
            >
                <AlertDialogOverlay>
                    <AlertDialogContent>
                        <AlertDialogHeader fontSize="lg" fontWeight="bold">
                            Delete Comment
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