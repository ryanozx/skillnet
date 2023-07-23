import {
    Box,
    VStack,
} from '@chakra-ui/react';
import InfoSection from './InfoSection';
import ProjectDisplay from '../ProjectDisplay/ProjectDisplay';
import { useEffect, useState } from 'react';
import axios from 'axios';
import { User } from '../../types';


export default function ProfileInfo({username}: {username: string}) {
   
    const [user, setUser] = useState<User>({
        AboutMe: "AboutMe not available",
        Email: "",
        Name: "Anonymous User",
        Title: "Title not available",
        ProfilePic: "",
        Username: "",
        ShowAboutMe: false,
        ShowTitle: false,
    });
    const [profileState, setProfileState] = useState<string>("loading");
    const [loadedUser, setLoadedUser] = useState<boolean>(false);

    useEffect(() => {
        const base_url = process.env.BACKEND_BASE_URL;
        const currentUrl = base_url + "/auth/user";
        if (username && username !== "undefined") {
            axios.get(currentUrl, { withCredentials: true })
            .then((res) => {
                const currentUser = res.data.data.Username;
                const isMyProfile = currentUser === username;
                const profileUrl = base_url + "/auth/users/" + username;
                axios.get(profileUrl, {withCredentials: true}).then((res) => {
                    setUser({...res.data.data, 
                        AboutMe: (res.data.data["ShowAboutMe"] || isMyProfile)? res.data.data["AboutMe"] : "",
                        Title: (res.data.data["ShowTitle"] || isMyProfile)? res.data.data["Title"] : "",
                    });
                    // Compare profile user to current user
                    
                    if (isMyProfile) {
                        setProfileState("self");
                    } else {
                        setProfileState("other");
                    }
                })
                .then(() => setLoadedUser(true))
                .catch((err) => {    
                    setProfileState("invalid")
                    console.log("booo")
                    console.log(err)
                });
            })
            .catch((err) => {
                console.log(err);
            })
        }
        // Fetch current user
    }, [username]); 

    if (profileState === "invalid") {
        return (
            <p>invalid username</p>
        )
    }

    return (
        <Box mt={10} mx={5} p={4}>
          <VStack spacing={10} align="start">
            <InfoSection
              user={user}
              username={username}
              {...(profileState === "self" && { setUser })}
            />
          </VStack>
          {loadedUser && <ProjectDisplay username={username} />}
        </Box>
      );
};
