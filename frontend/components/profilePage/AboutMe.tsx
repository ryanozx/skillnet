import { useState } from 'react';
import {
  Box,
  VStack,
  Text,
  Heading,
  Button
} from '@chakra-ui/react';

export default function AboutMe({ user }: any) {
    const [showMore, setShowMore] = useState(false);
    const [height, setHeight] = useState("200px");
    const { about = 'Description not available' } = user || {};

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
                    <Text>{about}</Text>
                </Box>
                {about && about.length > 100 && (
                    <Button onClick={handleClick} size="sm" alignSelf="flex-end">
                        {showMore ? "Show less" : "Show more"}
                    </Button>
                )}
            </VStack>
        </Box>
    );
}
