import { useState } from 'react';
import { Box, Heading, List, ListItem, Button } from '@chakra-ui/react';

export default function FollowedCommunities(props) {
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
            <List spacing={2} px={4}>
                {displayedCommunities.map((community, index) => (
                <ListItem key={index}>{community}</ListItem>
                ))}
            </List>
            {followedCommunities.length > 5 && (
                <Button onClick={handleClick} size="sm" alignSelf="flex-end">
                    {showMore ? "Show less" : "Show more"}
                </Button>
            )}
        </Box>
  );
}
