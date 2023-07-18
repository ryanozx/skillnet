import React from "react";
import {AlertDialog, AlertDialogBody, AlertDialogFooter, AlertDialogHeader, AlertDialogOverlay, Button, useDisclosure} from "@chakra-ui/react";
import {CloseIcon} from "@chakra-ui/icons";

interface EditPostCancelProps {
    handleClose: () => void;
}

export default function EditPostCancel(props : EditPostCancelProps) {
    return (
        <Button onClick={props.handleClose} colorScheme="red" leftIcon={<CloseIcon />}>
            Discard Changes
        </Button>
        )
}