import React, {useState, useEffect} from "react";
import {
    Box,
    Button,
    Card,
    CardHeader,
    CardBody,
    CardFooter,
    Divider,
    Flex,
    Heading,
    Spacer,
    Text,
} from "@chakra-ui/react";
import axios from "axios";
import EditProjectInfoModal from "./EditProjectInfoModal";
import DeleteProjectButton from "./DeleteProjectButton";

interface ProjectInfoProps {
    ProjectID: number,
    setProjectLoaded: React.Dispatch<React.SetStateAction<boolean>>
    setCommunityID: React.Dispatch<React.SetStateAction<number>>
}

export interface Project {
    ID: number
    Name: string
    ProjectInfo: string
}

export default function ProjectInfo(props : ProjectInfoProps) {
    const [project, setProject] = useState<Project>({
        ID: 0,
        Name: "Error",
        ProjectInfo: "There doesn't seem to be a project here...",
    });
    const [isOwner, setIsOwner] = useState<boolean>(false);
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [editModalOpen, setEditModalOpen] = useState<boolean>(false);

    const loadProject = async() => {
        if (!isLoading) {
            setIsLoading(true);
            const base_url = process.env.BACKEND_BASE_URL;
            const projectURL = base_url + "/auth/projects/" + props.ProjectID;
            const fetchData = axios.get(projectURL, {withCredentials: true});
            fetchData
            .then(res => {
                console.log(res.data.data)
                setProject({...res.data.data});
                setIsOwner(res.data.data["IsOwner"]);
                props.setCommunityID(res.data.data["CommunityID"]);
                props.setProjectLoaded(true);
            })
            .catch((error) => {
                console.log(error);
            })
            .finally(() => setIsLoading(false))
        }
    }
    

    useEffect(() => {
        loadProject()
    }, [props.ProjectID]);

    return (<Box w="100%" px={10} paddingBlockStart={5}>
    <Card>
        <CardHeader>
            <Flex>
                <Heading size="lg">{project.Name}</Heading>
                <Spacer />
                {isOwner && <>
                    <Button onClick={() => setEditModalOpen(true)}>Edit Project</Button>
                    <DeleteProjectButton projectID={props.ProjectID}/>
                    <EditProjectInfoModal isOpen={editModalOpen} setIsOpen={setEditModalOpen} project={project} updateProjectHandler={setProject}/>
                </>}
            </Flex>
        </CardHeader>
        <CardBody>
            <Divider />
            <Heading size="md" paddingBlock={3}>About Project</Heading>
            <Text style={{overflowWrap: "anywhere"}}>{project.ProjectInfo}</Text>
        </CardBody>
        {isOwner && <CardFooter>
            
        </CardFooter>}
    </Card>
</Box>)
}