import { useState } from 'react';
import { 
    Box, 
    Flex,
    Heading, 
    List, 
    ListItem, 
    Button,
    Link, 
    Divider
} from '@chakra-ui/react';

export default function FollowedCommunitiesList(props) {
    const [showMore, setShowMore] = useState(false);
    const handleClick = () => setShowMore(!showMore);

    const followedCommunities = [
        "Gardening",
        "Cooking",
        "Gaming",
        "Music",
        "Art",
        "Photography",
        "Sports",
        "Movies"
    ]

    const displayedCommunities = showMore ? followedCommunities : followedCommunities.slice(0, 5);

    return (
        <Box>
            <Heading size="md">Followed Communities</Heading>
            <List spacing={2} p={4}>
                {displayedCommunities.map((community, index) => (
                <ListItem key={index} py={1}>
                    <Link>{community}</Link>
                </ListItem>
                ))}
            </List>
            
            {followedCommunities.length > 5 && (
                <>
                    <Divider/>
                    <Flex justifyContent={"flex-end"} pt={3}>
                        <Button onClick={handleClick} size="sm" alignSelf={"flex-end"}>
                            {showMore ? "Show less" : "Show more"}
                        </Button>
                    </Flex>
                </>
            )}
        </Box>
  );
}
