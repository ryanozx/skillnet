import { Divider, Heading, Text } from '@chakra-ui/react';

export const HomePageHeader = () => (
    <Heading
        fontWeight={600}
        fontSize={{ base: '2xl', sm: '4xl', md: '6xl' }}
        lineHeight={'110%'}>
            SkillNet <br />
        <Divider colorScheme='black' py={2}/>
        <Text as={'span'} color={'green.400'}>
            Empower your ideas
        </Text>
    </Heading>
);