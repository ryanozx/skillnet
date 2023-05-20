import { 
    Box, 
    Button, 
    Grid,
    Flex
 } from "@chakra-ui/react";
import ProjectDisplayCard from "./ProjectDisplayCard";

export default function ProjectDisplay ({ projects }) {
  return (
    <Box>
        <Grid 
            templateColumns={{ base: 'repeat(1, 1fr)', sm: 'repeat(1, 1fr)', md: 'repeat(2, 1fr)', lg: 'repeat(4, 1fr)'}}
            gap={6} 
            mb={4}
        >
            {/* {projects.map((project, index) => (
            <Box key={index} minW="200px" minH="300px">
                <ProjectDisplayCard
                logo={project.logo}
                name={project.name}
                category={project.category}
                backdrop={project.backdrop}
                />
            </Box>
            ))} */}
            <Box minW="300px" minH="400px">
                <ProjectDisplayCard/>
            </Box>
            <Box minW="300px" minH="400px">
                <ProjectDisplayCard/>
            </Box>
            <Box minW="300px" minH="400px">
                <ProjectDisplayCard/>
            </Box>
            <Box minW="300px" minH="400px">
                <ProjectDisplayCard/>
            </Box>

        </Grid>

        <Flex justifyContent="flex-end">
            <Button>See All</Button>
        </Flex>
      
    </Box>
  );
};
