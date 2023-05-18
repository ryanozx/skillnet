import React from 'react';
import {
  Box,
  VStack,
} from '@chakra-ui/react';
import EditProfileModal from './EditInfoModal';
import BasicInfo from './BasicInfo';
import AboutMe from './AboutMe';
import ProjectDisplay from './ProjectDisplay';

export default function ProfileInfo(user) {

  return (
    <Box mt={10} mx={5} p={4} >
        <VStack spacing={10} align="start">
            <BasicInfo user={user}></BasicInfo>
            <AboutMe user={user}></AboutMe>
            <EditProfileModal/>
            <ProjectDisplay></ProjectDisplay>
        </VStack>
        
    </Box>
  );
};
