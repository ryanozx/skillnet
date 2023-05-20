import React from "react";

import {
    Grid,
    GridItem
}   from '@chakra-ui/react';
import NavBar from "../NavBar/NavBar"
import SideBar from "../SideBar/SideBar";

export default function LayoutContainer({children}) {
    return (
        <Grid
            templateAreas={`"header header"
                            "nav main"
                            `}
            gridTemplateRows={'7vh 1fr'}
            gridTemplateColumns={'20vw 3fr'}
            h='100vh'
            gap='1'
            color='blackAlpha.700'
        >
            <GridItem zIndex={2} bg='orange.300' area={'header'}>
                <NavBar></NavBar>
            </GridItem>
            <GridItem zIndex={1} bg='pink.300' area={'nav'}>
                <SideBar></SideBar>
            </GridItem>
            <GridItem zIndex={1} bg='green.300' area={'main'}>
                {children}
            </GridItem>
        
        </Grid>  
    );

}