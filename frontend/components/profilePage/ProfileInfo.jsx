import {
  Box,
  VStack,
} from '@chakra-ui/react';
import BasicInfo from './BasicInfo';
import AboutMe from './AboutMe';
import ProjectDisplay from './ProjectDisplay';
import { useEffect, useState } from 'react';
import axios from 'axios';

export default function ProfileInfo({username}) {
    const [user, setUser] = useState(null);
    
    useEffect(() => {
        const url = ''
        console.log('API call to get user information given username');
        const session_id = sessionStorage.getItem('session_id');
        const fetchData = axios.post(url, {
          username: username,
          session_id: session_id,
        });
        fetchData
            .then(response => {
                setUser(response.data);
            })
            .catch(error => {
                console.error(error);
            });
      }, [username]);

    if (!user) {
        return <div>Invalid User Id</div>;
    }

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
