import React from "react";
import {
    Grid,
    GridItem,
    useBreakpointValue
}   from '@chakra-ui/react';
import NavBar from "../NavBar/NavBar"
import SideBar from "../SideBar/SideBar";

export default function LayoutContainer({children}) {
    const templateColumns = useBreakpointValue({ base: '1fr', lg: '20vw 3fr' });
    const templateAreas = useBreakpointValue({ base: `"header" "main"`, lg: `"header header" "nav main"` });

    return (
        <Grid
            templateAreas={templateAreas}
            gridTemplateRows='60px 1fr'
            gridTemplateColumns={templateColumns}
            gap='1'
            color='blackAlpha.700'
        >
            <GridItem zIndex={2} bg='orange.300' area='header'>
                <NavBar />
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