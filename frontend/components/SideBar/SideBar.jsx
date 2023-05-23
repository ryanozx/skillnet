import { 
    Box, 
    VStack, 
    Heading, 
    List, 
    ListItem, 
    Link,
} from "@chakra-ui/react";
import Searchbar from '../base/Searchbar';
import FollowedCommunities from "./FollowedCommunities";
import CreateCommunityBtn from "./CreateCommunityBtn";
import { useEffect, useState } from "react";

export default function SideBar() {
    const [popularCommunities, setPopularCommunities] = useState([]);

    useEffect(() => {
        axios.get('/api/popular-communities')
        .then(response => {
            setPopularCommunities(response.data);
        })
        .catch(error => {
            console.error(error);
        });
    }, []);

    // followed communities
    // popular communities

    return (
        <Box p={5}>
            <VStack align="stretch" spacing={4}>
                <Heading size="md">Search Communities</Heading>
                <Searchbar/>

                <Heading size="md">Popular Communities</Heading>
                <List spacing={2} px={4}>
                    {/* {popularCommunities.map((community, index) => (
                        <ListItem key={index}>{community}</ListItem>
                    ))} */}
                    <ListItem><Link href="#">r/programming</Link></ListItem>
                </List>
                <FollowedCommunities/>
                <CreateCommunityBtn/>                
            </VStack>
        </Box>
  );
}
