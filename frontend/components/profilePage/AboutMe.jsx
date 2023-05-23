import { useState } from 'react';
import {
  Box,
  VStack,
  Text,
  Heading,
  Button
} from '@chakra-ui/react';

export default function AboutMe({ user }) {
    const [showMore, setShowMore] = useState(false);
    const [height, setHeight] = useState("200px");

    // const { about = `
    //     Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature 
    //     from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, 
    //     looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of 
    //     the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 
    //     of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise 
    //     on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", 
    //     comes from a line in section 1.10.32.
    //     The standard chunk of Lorem Ipsum used since the 1500s is reproduced below for those interested. Sections 1.10.32 and 
    //     1.10.33 from "de Finibus Bonorum et Malorum" by Cicero are also reproduced in their exact original form, accompanied by 
    //     English versions from the 1914 translation by H. Rackham.
    //     The standard chunk of Lorem Ipsum used since the 1500s is reproduced below for those interested. Sections 1.10.32 and 
    //     1.10.33 from "de Finibus Bonorum et Malorum" by Cicero are also reproduced in their exact original form, accompanied by 
    //     English versions from the 1914 translation by H. Rackham.
    //     The standard chunk of Lorem Ipsum used since the 1500s is reproduced below for those interested. Sections 1.10.32 and 
    //     1.10.33 from "de Finibus Bonorum et Malorum" by Cicero are also reproduced in their exact original form, accompanied by 
    //     English versions from the 1914 translation by H. Rackham.
    // `} = user || {};
    const { about = "this is a shorter description" } = user || {};

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
