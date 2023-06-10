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

export default function ProjectDisplayCard(project: any) {
    const { logo, name, category, backdrop } = project;
    return (
        <Card maxW='300px' h="400px">
            <CardBody>
                <Image
                    src={logo}
                    alt='Project logo'
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
    );
};
