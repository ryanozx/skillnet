import React from 'react';
import {
    Box,
    Container,
    Stack,
} from '@chakra-ui/react';
import CallToActionButtons from './CallToActionButtons';
import HomePageHeader from './HomePageHeader';
import { preventAuthAccess } from '../../withAuthRedirect';

export default preventAuthAccess(function HomePageContainer() {

    return (
        <>
            <Container maxW={'3xl'}>
                <Stack
                    as={Box}
                    textAlign={'center'}
                    spacing={{ base: 8, md: 14 }}
                    py={{ base: 20, md: 36 }}>
                    <HomePageHeader />
                    <CallToActionButtons />        
                </Stack>
            </Container>
        </>
    );
});

