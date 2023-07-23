import React, { useState } from "react";
import {
    Box,
    Button,
} from "@chakra-ui/react";
import { AddIcon } from "@chakra-ui/icons";
import CreateCommunityModal from "../../communityPage/CreateCommunityModal";

export default function CreateCommunityBtn() {
    const [isOpen, setIsOpen] = useState<boolean>(false);     
    const handleOpen = () => setIsOpen(true);

    return (
        <Box>
            <Button shadow='md' w='100%' leftIcon={<AddIcon />} colorScheme="green" onClick={handleOpen}>
                Create Community
            </Button>
            <CreateCommunityModal isOpen={isOpen} setIsOpen={setIsOpen}/>
        </Box>
    );
}