import React, { useState, useEffect } from "react";
import {
    Box,
    Card,
    CardBody,
    CardHeader,
    Grid,
    Flex,
    Heading,
    Button,
    Spacer,
    useDisclosure
} from "@chakra-ui/react";
import ProjectDisplayCard from './ProjectDisplayCard';
import ProjectDisplayModal from './ProjectDisplayModal';
import { ProjectMinimal } from "../../types";
import {AddIcon} from "@chakra-ui/icons";
import axios from "axios";
import CreateProjectModal from "./CreateProjectModal";

interface ProjectDisplayProps {
    communityID?: number
    username?: string
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
    
    const [projects, setProjects] = useState<ProjectMinimal[]>([]);
    const [displayProjects, setDisplayProjects] = useState<ProjectMinimal[]>(projects.slice(0, 4));
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [createModalOpen, setCreateModalOpen] = useState<boolean>(false);
    const [initialUrlUpdated, setInitialUrlUpdated] = useState<boolean>(false);

    const handleSeeAllClick = () => {
        onOpen();
    }

    const base_url = process.env.BACKEND_BASE_URL;
    const [url, setURL] = useState(base_url + "/auth/projects")

    useEffect(() => {
        if (props.communityID && props.communityID != 0) {
            setURL(url + "?community=" + props.communityID)
        }
        else if (props.username) {
            setURL(url + "?username=" + props.username)
        }
        setInitialUrlUpdated(true);
    }, [props.communityID, props.username])
    
    useEffect(() => {
        if (initialUrlUpdated) {
            updateProjects();
        }
    }, [initialUrlUpdated]);

    const updateProjects = async() => {
        if (!isLoading) {
            setIsLoading(true);

            const fetchData = axios.get(url, {withCredentials: true});
            fetchData
            .then((response) => {
                console.log(response.data.data)
                if (response.data["data"]["projects"] != null)
                {
                    setProjects([...projects, ...response.data["data"]["projects"]]);
                    
                }
                setURL(response.data["data"]["NextPageURL"]);
            })
            .catch((error) => {
                console.log(error);
            })
            .finally(() =>
                setIsLoading(false)
            )
        }
    }

    useEffect(() => {
        setDisplayProjects(projects.slice(0, 4));
    }, [projects]);
    
    return (
        <Card marginBlock={10}>
            <CardHeader>
                <Flex>
                    <Heading size="lg">Projects</Heading>
                    {props.communityID &&
                        <>
                            <Spacer />
                            <Button
                                leftIcon={<AddIcon></AddIcon>}
                                onClick={() => setCreateModalOpen(true)}
                            >Create Project</Button>
                            {createModalOpen && <CreateProjectModal isOpen={createModalOpen} setIsOpen={setCreateModalOpen} communityID={props.communityID}/>}
                        </>}
                </Flex>
            </CardHeader>
            <CardBody>
                {projects.length === 0 ? "There are no projects so far..." : <Grid 
                    templateColumns={{ base: 'repeat(2, 1fr)', md:'repeat(3, 1fr)', lg: 'repeat(4, 1fr)'}}
                    gap={6} 
                    mb={4} 
                >
                    {displayProjects.map((project: ProjectMinimal) => (
                        <Box 
                            key={project.ID}
                            width={{ base: "100%", sm: "auto" }}
                            minWidth={{ base: "250px", sm: "auto" }}
                        >
                            <ProjectDisplayCard
                                {...project}
                            />
                        </Box>
                    ))}
                </Grid>}

                {projects.length > 4 && 
                    <Flex justifyContent="flex-end">
                        <Button onClick={handleSeeAllClick}>See All</Button>
                    </Flex>
                }

                <ProjectDisplayModal
                    isOpen={isOpen}
                    onClose={onClose}
                    projects={projects}
                    updateProjects={updateProjects}
                />
            </CardBody>
        </Card>);
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


