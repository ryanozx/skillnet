import React from "react";
import { 
    Box, 
    Heading, 
    List, 
    ListItem, 
    Link,
} from "@chakra-ui/react";
import { useEffect, useState } from "react";
import axios from 'axios';

export default function PopularCommunitiesList() {

    const [popularCommunities, setPopularCommunities] = useState([]);

    

    // useEffect(() => {
    //     console.log('API call to get list of popular communities');
    //     const url = '';
    //     axios.get('/api/popular-communities')
    //     .then(response => {
    //         setPopularCommunities(response.data);
    //     })
    //     .catch(error => {
    //         console.error(error);
    //     });
    // }, []);

    return (
        <>
            <Heading size="md">Popular Communities</Heading>
            <List spacing={2} px={4}>
                {/* {popularCommunities.map((community, index) => (
                    <ListItem key={index}>{community}</ListItem>
                ))} */}
                <ListItem><Link href="#">r/programming</Link></ListItem>
            </List>
        </>
    )
}