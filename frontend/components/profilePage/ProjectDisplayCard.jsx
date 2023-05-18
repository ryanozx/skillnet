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

const ProjectDisplayCard = ( project ) => {
    
    const name = "SkillNet"
    const category = "Web Development"
    const logo = "https://via.placeholder.com/150x150/000000/FFFFFF?text=Logo"
    const backdrop= "" 
    


  return (
    <Card maxW='300px' h="400px">
        <CardBody>
            <Image
                src='https://images.unsplash.com/photo-1555041469-a586c61ea9bc?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1770&q=80'
                alt='Green double couch with wooden legs'
                borderRadius='lg'
                h="250px"
                w="100%"
                objectFit={'cover'}
            />
            <Stack mt='6' spacing='1'>
                <Heading size='md'>{name}</Heading>
                <Divider/>
                <Text>
                    {category}
                </Text>
                
            </Stack>
            <Flex
                float="right"
            >
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
  );
};

export default ProjectDisplayCard;
