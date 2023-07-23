import { useState } from 'react';
import {
  Box,
  VStack,
  Text,
  Heading,
  Button
} from '@chakra-ui/react';

export default function AboutMe({ aboutMe }: { aboutMe: string }) {
    const [showMore, setShowMore] = useState(false);
    const [height, setHeight] = useState("200px");
    const handleClick = () => {
        setShowMore(!showMore);
        setHeight(showMore ? "200px" : "auto");
    };

    return (
        <Box
          w="100%"
          p={4}
          mb={4}
        >
            <VStack spacing={5} align="start">
                <Heading size="md" px={2}>About Me</Heading>
                <Box 
                    bg="green.200"
                    w="100%"
                    p={5}
                    h={height}
                    overflow="hidden"
                >
                    <Text>{aboutMe ? aboutMe : "No description available"}</Text>
                </Box>
                {aboutMe && aboutMe.length > 100 && (
                    <Button onClick={handleClick} size="sm" alignSelf="flex-end">
                        {showMore ? "Show less" : "Show more"}
                    </Button>
                )}
            </VStack>
        </Box>
    );
}
