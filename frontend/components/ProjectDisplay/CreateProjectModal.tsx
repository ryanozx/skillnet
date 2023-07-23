import React, {useState} from "react";
import {
    Button,
    FormControl,
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
import { escapeHtml } from "../../types";

interface CreateProjectModalProps {
    isOpen: boolean,
    setIsOpen: React.Dispatch<React.SetStateAction<boolean>>
    communityID: number,
}

export default function CreateProjectModal(props : CreateProjectModalProps) {
    const toast = useToast();
    const router = useRouter();

    const [projectName, setProjectName] = useState<string>("");
    const [aboutProject, setAboutProject] = useState<string>("");
    const baseURL = process.env.BACKEND_BASE_URL;
    const createProjectURL = baseURL + "/auth/projects"; 

    const handleClose = () => {
        setProjectName("");
        setAboutProject("");
        props.setIsOpen(false);
    };

    const handleProjectNameChange = (e : React.ChangeEvent<HTMLInputElement>) => {
        setProjectName(e.target.value);
    }

    const handleCreate = () => {
        axios.post(createProjectURL, { 
                "Name": escapeHtml(projectName),
                "About": escapeHtml(aboutProject),
                "communityID": props.communityID }, {withCredentials: true})
            .then(res => {
                toast({
                title: "Community created.",
                description: "We've created your community for you.",
                status: "success",
                duration: 5000,
                isClosable: true,
                });
                console.log(res.data.data);
                router.push("/projects/" + res.data.data["ID"]);
                handleClose();
            })
            .catch((error) => {
                console.error(error);
                handleClose();
                toast({
                    title: "Project not created.",
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
                    <ModalHeader>Create a new Project</ModalHeader>
                    <ModalCloseButton />
                    <ModalBody>
                        <FormControl isInvalid={projectName !== ""}>
                            <FormLabel>Project Name</FormLabel>
                            <Input placeholder="Enter project name" value={projectName} onChange={handleProjectNameChange} />
                        </FormControl>
                        <FormControl mt={4}>
                            <FormLabel>About Project</FormLabel>
                            <Textarea placeholder="Enter About Project" value={aboutProject} onChange={(e) => {setAboutProject(e.target.value)}} />
                        </FormControl>
                    </ModalBody>

                    <ModalFooter>
                        <Button colorScheme="blue" mr={3} onClick={handleCreate} isDisabled={projectName === ""}>
                            Create
                        </Button>
                        <Button variant="ghost" onClick={handleClose}>Cancel</Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
    )
}
