import {
    Box,
    Button
} from "@chakra-ui/react";
import { AddIcon } from "@chakra-ui/icons";

export default function CreateCommunityBtn() {
    return (
        <Button leftIcon={<AddIcon />} colorScheme="green">
            Create Community
        </Button>
    );
}