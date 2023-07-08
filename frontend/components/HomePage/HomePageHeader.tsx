import React from 'react';
import { Divider, Heading, Text } from '@chakra-ui/react';

const HomePageHeader = () => (
    <Heading 
        data-testid="home-page-header"
        fontWeight={600}
        fontSize={{ base: '2xl', sm: '4xl', md: '6xl' }}
        lineHeight={'110%'}>
            SkillNet <br />
        <Divider colorScheme='black' py={2}/>
        <Text data-testid="home-page-subheader" as={'span'} color={'green.400'}>
            Empower your ideas
        </Text>
    </Heading>
);

export default HomePageHeader;