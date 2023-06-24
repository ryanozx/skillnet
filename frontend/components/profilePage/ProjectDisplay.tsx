import React, { useState } from "react";
import {
  Box,
  Grid,
  Flex,
  Button,
  useDisclosure
} from "@chakra-ui/react";
import ProjectDisplayCard from './ProjectDisplayCard';
import ProjectDisplayModal from './ProjectDisplayModal';

interface ProjectDisplayProps {
    projects: ProjectView[];  // Make projects optional
}

interface ProjectView {
    logo: string;
    name: string;
    category: string;
}

export default function ProjectDisplay (props: ProjectDisplayProps) {    
    const { isOpen, onOpen, onClose } = useDisclosure();
    // const projects = [ 
    //         {
    //             logo: "https://images.unsplash.com/photo-1555041469-a586c61ea9bc?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1770&q=80",
    //             name: "SkillNet", 
    //             category: "Web Development", 
    //         },
    //         {
    //             logo: "https://images.unsplash.com/photo-1555041469-a586c61ea9bc?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1770&q=80",
    //             name: "SkillNet 2", 
    //             category: "Pencil Art", 
    //         },
    //         {
    //             logo: "https://images.unsplash.com/photo-1555041469-a586c61ea9bc?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1770&q=80",
    //             name: "SkillNet 3", 
    //             category: "Gardening", 
    //         },
    //         {
    //             logo: "https://images.unsplash.com/photo-1555041469-a586c61ea9bc?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1770&q=80",
    //             name: "SkillNet 4", 
    //             category: "Cooking", 
    //         },
    //         {
    //             logo: "https://images.unsplash.com/photo-1555041469-a586c61ea9bc?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1770&q=80",
    //             name: "SkillNet 5", 
    //             category: "Web Development", 
    //         },
    //     ]
    
    const projects = props.projects || ([] as ProjectView[]);
    const [displayProjects, setDisplayProjects] = useState<ProjectView[]>(projects.slice(0, 4));

    const handleSeeAllClick = () => {
        onOpen();
    }
    
    return (
        <Box>
            
            <Grid 
                templateColumns={{ base: 'repeat(2, 1fr)', md:'repeat(3, 1fr)', lg: 'repeat(4, 1fr)'}}
                gap={6} 
                mb={4} 
            >
                {displayProjects.map((project: any, index: any) => (
                    <Box 
                        key={index}
                        width={{ base: "100%", sm: "auto" }}
                        minWidth={{ base: "250px", sm: "auto" }}
                    >
                        <ProjectDisplayCard
                            logo={project.logo}
                            name={project.name}
                            category={project.category}
                            backdrop={project.backdrop}
                        />
                    </Box>
                ))}
            </Grid>

            {projects.length > 4 && 
                <Flex justifyContent="flex-end">
                    <Button onClick={handleSeeAllClick}>See All</Button>
                </Flex>
            }

            <ProjectDisplayModal
                isOpen={isOpen}
                onClose={onClose}
                projects={projects}
            />
        </Box>
    );
};

// const projects = [ 
    //         {
    //             logo: "https://images.unsplash.com/photo-1555041469-a586c61ea9bc?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1770&q=80",
    //             name: "SkillNet", 
    //             category: "Web Development", 
    //         },
    //         {
    //             logo: "https://images.unsplash.com/photo-1555041469-a586c61ea9bc?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1770&q=80",
    //             name: "SkillNet 2", 
    //             category: "Pencil Art", 
    //         },
    //         {
    //             logo: "https://images.unsplash.com/photo-1555041469-a586c61ea9bc?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1770&q=80",
    //             name: "SkillNet 3", 
    //             category: "Gardening", 
    //         },
    //         {
    //             logo: "https://images.unsplash.com/photo-1555041469-a586c61ea9bc?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1770&q=80",
    //             name: "SkillNet 4", 
    //             category: "Cooking", 
    //         },
    //         {
    //             logo: "https://images.unsplash.com/photo-1555041469-a586c61ea9bc?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1770&q=80",
    //             name: "SkillNet 5", 
    //             category: "Web Development", 
    //         },
    //     ]


