import React from "react";
import { 
    Box, 
    VStack, 
    Heading, 
} from "@chakra-ui/react";
import PopularCommunitiesList from "./PopularCommunitiesList";
import CreateCommunityBtn from "./CreateCommunityBtn";

export default function SideBar() {

    return (
        <Box p={5}>
            <VStack align="stretch" spacing={4}>
                <PopularCommunitiesList/>
                <CreateCommunityBtn/>                
            </VStack>
        </Box>
  );
}
