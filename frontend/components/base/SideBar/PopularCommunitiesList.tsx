import React from "react";
import { 
    Heading, 
    List, 
    ListItem, 
    Link,
} from "@chakra-ui/react";
import { useEffect, useState } from "react";
import axios from 'axios';
import { Community } from "../../communityPage/CommunityInfo";

interface CommunityView {
    Community: Community,
    IsOwner: boolean,
}

export default function PopularCommunitiesList() {
    const [popularCommunities, setPopularCommunities] = useState<React.JSX.Element[]>([]);
    const baseURL = process.env.BACKEND_BASE_URL;
    const [url, setURL] = useState<string>(baseURL + "/auth/community")
    useEffect(() => {
        axios.get(url, {withCredentials: true})
        .then(res => {
            setPopularCommunities([...popularCommunities, ...res.data.data["Communities"].map(
                (community : CommunityView) => 
                    <ListItem key={community.Community.ID}>
                        <Link href={`/communities/${community.Community.Name}`}>{community.Community.Name}</Link> 
                    </ListItem>)
            ]);
        })
        .catch(error => {
            console.error(error);
        });
    }, []);

    return (
        <>
            <Heading size="md">Recent Communities</Heading>
            <List spacing={2} px={4}>
                {popularCommunities}
            </List>
        </>
    )
}