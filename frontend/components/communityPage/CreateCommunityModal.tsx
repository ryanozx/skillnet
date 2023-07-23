import React, {useState} from "react";
import {
    Button,
    FormControl,
    FormErrorMessage,
    FormLabel,
    Input,
    Modal,
    ModalBody,
    ModalCloseButton,
    ModalContent,
    ModalFooter,
    ModalHeader,
    ModalOverlay,
    Textarea,
    useToast
} from "@chakra-ui/react"
import { useRouter } from "next/router";
import axios from "axios";

interface CreateCommunityModalProps {
    isOpen: boolean,
    setIsOpen: React.Dispatch<React.SetStateAction<boolean>>
}

export default function CreateCommunityModal(props : CreateCommunityModalProps) {
    const toast = useToast();
    const router = useRouter();

    const [communityName, setCommunityName] = useState<string>("");
    const [aboutCommunity, setAboutCommunity] = useState<string>("");
    const [validCommunityName, setValidCommunityName] = useState<boolean>(false);
    const baseURL = process.env.BACKEND_BASE_URL;
    const createCommunityURL = baseURL + "/auth/community"; 

    const handleClose = () => {
        setCommunityName("");
        setAboutCommunity("");
        props.setIsOpen(false);
    };

    const handleCommunityNameChange = (e : React.ChangeEvent<HTMLInputElement>) => {
        setCommunityName(e.target.value);
        setValidCommunityName(validateCommunityName(e.target.value));
    }

    const handleCreate = () => {
        axios.post(createCommunityURL, { 
                "Name": communityName,
                "About": aboutCommunity }, {withCredentials: true})
            .then(() => {
                toast({
                title: "Community created.",
                description: "We've created your community for you.",
                status: "success",
                duration: 5000,
                isClosable: true,
                });
                router.push("/communities/" + communityName);
                handleClose();
            })
            .catch((error) => {
                console.error(error);
                handleClose();
                toast({
                    title: "Community not created.",
                    description: "we encountered an error: " + error.message + ".",
                    status: "error",
                    duration: 5000,
                    isClosable: true,
                });
            });
      };

    return (
        <Modal isOpen={props.isOpen} onClose={handleClose}>
                <ModalOverlay />
                <ModalContent>
                    <ModalHeader>Create a new Community</ModalHeader>
                    <ModalCloseButton />
                    <ModalBody>
                        <FormControl isInvalid={communityName !== "" && !validCommunityName}>
                            <FormLabel>Community Name</FormLabel>
                            <Input placeholder="Enter community name" value={communityName} onChange={handleCommunityNameChange} />
                            {!validCommunityName &&
                                <FormErrorMessage>Name can only contain alphabetical characters.</FormErrorMessage>
                            }
                        </FormControl>
                        <FormControl mt={4}>
                            <FormLabel>About Community</FormLabel>
                            <Textarea placeholder="Enter About Community" value={aboutCommunity} onChange={(e) => {setAboutCommunity(e.target.value)}} />
                        </FormControl>
                    </ModalBody>

                    <ModalFooter>
                        <Button colorScheme="blue" mr={3} onClick={handleCreate} isDisabled={!validCommunityName}>
                            Create
                        </Button>
                        <Button variant="ghost" onClick={handleClose}>Cancel</Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
    )
}

function validateCommunityName(name : string) {
    const regex = /^[A-Za-z]+$/i;
    return regex.test(name)
}