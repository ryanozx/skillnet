import {
    Grid,
    GridItem,
    useBreakpointValue
}   from '@chakra-ui/react';
import NavBar from "../NavBar/NavBar"
import SideBar from "../SideBar/SideBar";
import React, { useEffect, ReactNode } from 'react';
// import { useSelector, useDispatch } from 'react-redux';
import axios from 'axios';
// import { loginRequest, loginSuccess, loginFailure } from '../../actions/userActions';
// import { RootState } from '../../reducers/rootReducer';

interface DefaultLayoutContainerProps {
    children: ReactNode;
}

export default function DefaultLayoutContainer({ children }: DefaultLayoutContainerProps) {

    const templateColumns = useBreakpointValue({ base: '1fr', lg: '20vw 3fr' });
    const templateAreas = useBreakpointValue({ base: `"header" "main"`, lg: `"header header" "nav main"` });
    const isLoggedIn = false;
    const user = null;
    // const dispatch = useDispatch();
    // const userRedux = useSelector((state: RootState) => state.user);
    // const {loading, isLoggedIn, user, error} = userRedux;

    // useEffect(() => {
    //     dispatch(loginRequest());
    //     // url for session id validation
    //     console.log('API call to check if user is logged in');
    //     const url = '';
    //     const sessionId = sessionStorage.getItem('sessionId');
    //     if (sessionId) {
    //     axios.post(url, {}, {
    //         headers: {
    //         Authorization: `Bearer ${sessionId}`
    //         }
    //     })
    //     .then(response => {
    //         const user = response.data;
    //         dispatch(loginSuccess(user));
    //     })
    //     .catch(error => {
    //         dispatch(loginFailure(error.message));
    //         sessionStorage.removeItem('sessionId');
    //     });
    //     } else {
    //         dispatch(loginFailure(new Error('No session ID found')));
    //     }
    // }, []);

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
