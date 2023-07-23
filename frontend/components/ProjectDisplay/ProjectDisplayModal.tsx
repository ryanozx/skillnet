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
  Text,
} from "@chakra-ui/react";
import ProjectDisplayCard from './ProjectDisplayCard';
import { ProjectMinimal } from "../../types";
import InfiniteScroll from "react-infinite-scroll-component";

interface ProjectDisplayProps {
    isOpen: boolean,
    onClose: () => void,
    projects: ProjectMinimal[],
    updateProjects: () => Promise<void>,
    noMoreProjects: boolean,
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
                        <InfiniteScroll 
                            next={props.updateProjects} 
                            hasMore={!props.noMoreProjects} 
                            loader={
                                <Box paddingBlock="10px">
                                    <Text textAlign="center">Loading...</Text>
                                </Box>
                            }
                            dataLength={props.projects.length}
                            endMessage={
                                <Box paddingBlock="10px">
                                    <Text textAlign="center">No more posts to load.</Text>
                                </Box>}>

                        </InfiniteScroll>
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