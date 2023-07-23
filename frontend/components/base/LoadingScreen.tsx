import React from 'react';
import { Flex, Spinner, Text, useColorModeValue } from '@chakra-ui/react';

interface LoadingScreenProps {
  loadingText?: string;
}

export default function LoadingScreen(props: LoadingScreenProps) {
    const { loadingText } = props;
    const color = useColorModeValue("black", "white");
    return (
        <Flex
        flexDirection="column"
        width="100vw"
        height="100vh"
        justifyContent="center"
        alignItems="center"
        backgroundColor={useColorModeValue("gray.200", "gray.800")}
        >
        <Spinner
            thickness="4px"
            speed="0.65s"
            emptyColor="gray.200"
            color={color}
            size="xl"
        />
        <Text
            marginTop={4}
            fontSize="lg"
            fontWeight="medium"
            color={color}
        >
            {loadingText}
        </Text>
        </Flex>
    );
}
