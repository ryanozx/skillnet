import {
  Box,
  VStack,
} from '@chakra-ui/react';
import BasicInfo from './BasicInfo';
import AboutMe from './AboutMe';
import ProjectDisplay from './ProjectDisplay';
import { useEffect, useState } from 'react';
import axios from 'axios';

export default function ProfileInfo({user_id}) {
    const [user, setUser] = useState(null);

    useEffect(() => {
        const sessionID = sessionStorage.getItem('sessionID');
    
        const fetchData = axios.post('/api/profile', {
          user_id,
          sessionID,
        });
    
        fetchData
          .then(response => {
            setProfile(response.data);
          })
          .catch(error => {
            console.error(error);
          });
    
      }, [user_id]);

    // if (!user) {
    //     return <div>Loading...</div>;
    // }

    console.log("detected user id " + user_id)
    return (
        <Box mt={10} mx={5} p={4} >
            <VStack spacing={10} align="start">
                <BasicInfo user={user}></BasicInfo>
                <AboutMe user={user}></AboutMe>
                <ProjectDisplay></ProjectDisplay>
            </VStack>    
        </Box>
    );
};
