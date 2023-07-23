import { useState } from "react";
import {
  Box,
  Grid,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalCloseButton,
  ModalBody,
  ModalHeader,
} from "@chakra-ui/react";
import ProjectDisplayCard from './ProjectDisplayCard';
import { ProjectMinimal } from "../../types";

interface ProjectDisplayProps {
    isOpen: boolean,
    onClose: () => void,
    projects: ProjectMinimal[],
    updateProjects: () => Promise<void>,
}

export default function ProjectDisplay (props: ProjectDisplayProps) {
    return (
        <Modal 
            isOpen={props.isOpen} 
            onClose={props.onClose}
            size = {{base: "full", md:"6xl"}}
        >
            <ModalOverlay />
            <ModalContent>
                <ModalCloseButton />
                <ModalHeader>All projects</ModalHeader>
                <ModalBody>
                    <Grid 
                        templateColumns={{ sm: 'repeat(2, 1fr)', md: 'repeat(3, 1fr)', lg: 'repeat(4, 1fr)'}}
                        gap={6} 
                        mb={4}
                    >
                        {props.projects.map((project: ProjectMinimal) => (
                            <Box key={project.ID}>
                                <ProjectDisplayCard {...project}/>
                            </Box>
                        ))}
                    </Grid>
                </ModalBody>
            </ModalContent>
        </Modal>
    );
}