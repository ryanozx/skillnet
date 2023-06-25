import {
    Grid,
    GridItem,
    useBreakpointValue
}   from '@chakra-ui/react';
import NavBar from "./NavBar/NavBar"
import SideBar from "./SideBar/SideBar";
import React, { useState, useEffect, ReactNode } from 'react';
import axios from 'axios';
import { useUser } from '../../userContext';
import { requireAuth } from '../../withAuthRedirect';

interface DefaultLayoutContainerProps {
    children: ReactNode;
}



export default requireAuth(function DefaultLayoutContainer({ children }: DefaultLayoutContainerProps) {
    const { needUpdate, setNeedUpdate } = useUser();
    const [ profilePic, setProfilePic ] = useState("");
    const [ username, setUsername ] = useState("");
    const templateColumns = useBreakpointValue({ base: '1fr', lg: '20vw 3fr' });
    const templateAreas = useBreakpointValue({ base: `"header" "main"`, lg: `"header header" "nav main"` });

    useEffect(() => {
        if (needUpdate) {
            const base_url = process.env.BACKEND_BASE_URL;
            axios.get(base_url + '/auth/user', { withCredentials: true })
                .then((res) => {
                    const { ProfilePic, Username } = res.data.data;
                    setProfilePic(ProfilePic);
                    setUsername(Username);
                    setNeedUpdate(false); // reset the flag after the data is updated
                })
                .catch((err) => {
                    console.log(err);
                });
        }
    }, [needUpdate]);
    return (
        <Grid
            templateAreas={templateAreas}
            gridTemplateRows='60px 1fr'
            gridTemplateColumns={templateColumns}
            gap='1'
            color='blackAlpha.700'
            minHeight='100vh'
        >
            <GridItem zIndex={2} bg='orange.300' area='header'>
                <NavBar profilePic={profilePic} username={username}/>
            </GridItem>
            {templateColumns !== '1fr' && (
                <GridItem zIndex={1} bg='pink.300' area='nav'>
                    <SideBar />
                </GridItem>
            )}
            <GridItem zIndex={1} bg='green.300' area='main'>
                {children}
            </GridItem>
        </Grid>  
    );
});
