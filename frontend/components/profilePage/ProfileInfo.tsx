import {
    Text,
  Box,
  VStack,
} from '@chakra-ui/react';
import InfoSection from './InfoSection';
import ProjectDisplay from './ProjectDisplay';
import { useEffect, useState } from 'react';
import axios from 'axios';
import { User } from '../../types';

export default function ProfileInfo({username}: {username: string}) {
    const [loading, setLoading] = useState<boolean>(true);
    const [user, setUser] = useState<User>({
        AboutMe: "",
        Email: "",
        Name: "",
        Title: "",
        ProfilePic: "",
        Username: "",
        Projects: []
    });

    useEffect(() => {
        const url = "http://localhost:8080/auth/user"
        
        axios.get(url, {withCredentials: true })
            .then((res) => {
                const { AboutMe, Email, Name, Title, ProfilePic, Username, Projects } = res.data.data;
                setUser({
                    AboutMe: AboutMe ? AboutMe : "No description available",
                    Email: Email,
                    Name: Name ? Name : "No display name",
                    Title: Title ? Title : "No title available",
                    ProfilePic: ProfilePic,
                    Username: Username, 
                    Projects: Projects,
                  });                
            }).catch((err) => {
                console.log(err);
            })
        
    }, []);

    return (
        <Box mt={10} mx={5} p={4} >
            <VStack spacing={10} align="start">
                <InfoSection 
                    user = {user}
                    setUser = {setUser}
                />
                <ProjectDisplay/>
            </VStack>    
        </Box>
    );
};
