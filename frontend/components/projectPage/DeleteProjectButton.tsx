import React, {useRef} from "react"
import axios from "axios"
import {AlertDialog, AlertDialogBody, AlertDialogContent, AlertDialogFooter, AlertDialogHeader, AlertDialogOverlay, Button, MenuItem, useToast, useDisclosure} from "@chakra-ui/react"
import { useRouter } from "next/router"

interface DeleteProjectButtonProps {
    projectID: number,
}

export default function DeleteProjectButton(props : DeleteProjectButtonProps) {
    const {isOpen, onOpen, onClose} = useDisclosure();
    const cancelRef = useRef<HTMLButtonElement>(null);
    const base_url = process.env.BACKEND_BASE_URL;
    const projectURL = base_url + `/auth/projects/${props.projectID}`;
    const toast = useToast();
    const router = useRouter();

    const handleDelete = () => {
        axios.delete(projectURL, {withCredentials: true})
        .then(() => {
            toast({
                title: "Post deleted.",
                description: "Your post has been deleted.",
                status: "success",
                duration: 5000,
                isClosable: true,
            });
            router.push("/feed");
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
            <Button colorScheme="red" marginInlineStart={5}
                onClick={onOpen}
            >Delete project
            </Button>

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