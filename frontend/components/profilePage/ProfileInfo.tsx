import {
  Box,
  VStack,
} from '@chakra-ui/react';
import BasicInfo from './BasicInfo';
import AboutMe from './AboutMe';
import ProjectDisplay from './ProjectDisplay';
import { useEffect, useState } from 'react';
import axios from 'axios';

type ProfileInfoProps = {
    username?: string,
    ownProfile: boolean
}

export default function ProfileInfo({username, ownProfile} : ProfileInfoProps) {
    const [user, setUser] = useState(null);
    
    useEffect(() => {
        const url = ownProfile ? "http://localhost:8080/auth/user" : `http://localhost:8080/users/${username}`;
        console.log('API call to get user information given username');
        const fetchData = axios.get(url, {withCredentials: true});
        fetchData
            .then(response => {
                console.log(response.data);
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
