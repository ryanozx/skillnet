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


export default function ProjectDisplay (props: any) {
    const { isOpen, onClose, projects } = props;
    return (
        <Modal 
            isOpen={isOpen} 
            onClose={onClose}
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
                        {projects.map((project: any, index: any) => (
                            <Box key={index}>
                                <ProjectDisplayCard
                                    logo={project.logo}
                                    name={project.name}
                                    category={project.category}
                                    backdrop={project.backdrop}
                                />
                            </Box>
                        ))}
                    </Grid>
                </ModalBody>
            </ModalContent>
        </Modal>
    );
}