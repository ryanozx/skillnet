import { 
    Box, 
    VStack, 
    Heading, 
} from "@chakra-ui/react";
import Searchbar from '../base/Searchbar';
import FollowedCommunitiesList from "./FollowedCommunitiesList";
import PopularCommunitiesList from "./PopularCommunitiesList";
import CreateCommunityBtn from "./CreateCommunityBtn";

export default function SideBar() {

    return (
        <Box p={5}>
            <VStack align="stretch" spacing={4}>
                <Heading size="md">Search Communities</Heading>
                <Searchbar/>
                <PopularCommunitiesList/>
                <FollowedCommunitiesList/>
                <CreateCommunityBtn/>                
            </VStack>
        </Box>
  );
}
