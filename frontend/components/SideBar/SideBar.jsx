import { 
    Box, 
    VStack, 
    Input, 
    Button, 
    Heading, 
    List, 
    ListItem, 
    Link,
    IconButton } from "@chakra-ui/react";
import { AddIcon } from "@chakra-ui/icons";
import Searchbar from '../base/Searchbar';
import FollowedCommunities from "./FollowedCommunities";

export default function SideBar() {

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
                    <ListItem>
                        <Link>Gardening</Link>
                    </ListItem>

                </List>
                <FollowedCommunities/>
                
            </VStack>
        </Box>
  );
}
