import {  
    Image, 
    Text, 
    Card,
    CardBody,
    Stack,
    Heading,
    Divider, 
    Button, 
    Flex
} from "@chakra-ui/react";
import { ProjectMinimal } from "../../types";

export default function ProjectDisplayCard(project: ProjectMinimal) {
    return (
        <a href={project.URL}>
            <Card maxW='300px' h="400px" _hover={{marginTop: -2, marginBottom: 2}}>
            <CardBody>
                <Image
                    src="https://images.unsplash.com/photo-1555041469-a586c61ea9bc?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1770&q=80"
                    alt='Project logo'
                    borderRadius='lg'
                    h="250px"
                    w="100%"
                    objectFit={'cover'}
                />
                <Stack mt='6' spacing='1'>
                    <Heading size='md'>{project.Name}</Heading>
                    <Divider/>
                    <Text>
                        {project.Community}
                    </Text>
                    
                </Stack>
                <Flex float="right">
                    <Button
                        colorScheme='blue'
                        variant='outline'
                        size='sm'
                        alignSelf="end"
                        
                    >
                        View
                    </Button>
                </Flex>
            </CardBody>
        </Card>
        </a> 
    );
};
