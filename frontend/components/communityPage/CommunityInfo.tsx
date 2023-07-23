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
import ProjectDisplay from "../ProjectDisplay/ProjectDisplay";
import EditCommunityModal from "./EditCommunityInfoModal";

interface CommunityInfoProps {
    communityName: string
    setCommunityID: React.Dispatch<React.SetStateAction<number>>
    setCommunityLoaded: React.Dispatch<React.SetStateAction<boolean>>
}

export interface Community {
    ID: number
    Name: string
    About: string
}

export default function CommunityInfo(props : CommunityInfoProps) {
    const [community, setCommunity] = useState<Community>({
        ID: 0,
        Name: "Error",
        About: "There doesn't seem to be a community here...",
    });
    const [isOwner, setIsOwner] = useState<boolean>(false);
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [loadedCommunity, setLoadedCommunity] = useState<boolean>(false);
    const [editModalOpen, setEditModalOpen] = useState<boolean>(false);

    const loadCommunity = async() => {
        if (!isLoading && props.communityName) {
            setIsLoading(true);
            const base_url = process.env.BACKEND_BASE_URL;
            const communityURL = base_url + "/auth/community/" + props.communityName
            const fetchData = axios.get(communityURL, {withCredentials: true});
            fetchData
            .then(res => {
                setCommunity({...res.data.data["Community"]});
                setIsOwner(res.data.data["IsOwner"]);
                props.setCommunityID(res.data.data["Community"]["ID"])
            })
            .then(() => {
                setLoadedCommunity(true);
                props.setCommunityLoaded(true);
            })
            .catch((error) => {
                console.log(error);
            })
            .finally(() => setIsLoading(false))
        }
    }
    

    useEffect(() => {
        loadCommunity()
    }, [props.communityName]);

    return (<Box w="100%" px={10} paddingBlockStart={5}>
    <Card>
        <CardHeader>
            <Flex>
                <Heading size="lg">{community.Name}</Heading>
                <Spacer />
                {isOwner && <>
                    <Button onClick={() => setEditModalOpen(true)}>Edit Community</Button>
                    <EditCommunityModal isOpen={editModalOpen} setIsOpen={setEditModalOpen} community={community} updateCommunityHandler={setCommunity}/>
                </>}
            </Flex>
        </CardHeader>
        <CardBody>
            <Divider />
            <Heading size="md" paddingBlock={3}>About Community</Heading>
            <Text style={{overflowWrap: "anywhere"}}>{community.About}</Text>
        </CardBody>
        {isOwner && <CardFooter>
            
        </CardFooter>}
    </Card>
        {loadedCommunity && <ProjectDisplay communityID={community.ID}></ProjectDisplay>}
    
</Box>)
}