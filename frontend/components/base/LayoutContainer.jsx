import {
    Grid,
    GridItem,
    useBreakpointValue
}   from '@chakra-ui/react';
import NavBar from "../NavBar/NavBar"
import SideBar from "../SideBar/SideBar";
import React, { useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import axios from 'axios';
import { loginRequest, loginSuccess, loginFailure } from '../../actions/userActions';

export default function LayoutContainer({children}) {

    const templateColumns = useBreakpointValue({ base: '1fr', lg: '20vw 3fr' });
    const templateAreas = useBreakpointValue({ base: `"header" "main"`, lg: `"header header" "nav main"` });
    const dispatch = useDispatch();
    const userRedux = useSelector(state => state.user);
    const {loading, isLoggedIn, user, error} = userRedux;

    useEffect(() => {
        dispatch(loginRequest());
        const sessionId = sessionStorage.getItem('sessionId');
        if (sessionId) {
        axios.post('your-backend-url/api/auth', {}, {
            headers: {
            Authorization: `Bearer ${sessionId}`
            }
        })
        .then(response => {
            const userRedux = response.data;
            dispatch(loginSuccess(userRedux));
            // Here, you can also dispatch other actions to update 
            // the other parts of your state (followed communities, 
            // profile picture, notifications)
        })
        .catch(error => {
            dispatch(loginFailure(error.message));
            sessionStorage.removeItem('sessionId');
        });
        } else {
            dispatch(loginFailure('No session ID found'));
        }
    }, []);

    return (
        <Grid
            templateAreas={templateAreas}
            gridTemplateRows='60px 1fr'
            gridTemplateColumns={templateColumns}
            gap='1'
            color='blackAlpha.700'
        >
            <GridItem zIndex={2} bg='orange.300' area='header'>
                <NavBar user={user} isLoggedIn={isLoggedIn}/>
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
}
