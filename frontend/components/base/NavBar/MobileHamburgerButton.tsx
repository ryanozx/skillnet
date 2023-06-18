import {
    IconButton,
    useDisclosure,
} from '@chakra-ui/react';
import React from 'react';
import {
    HamburgerIcon,
    CloseIcon,
} from '@chakra-ui/icons';

interface HamburgerButtonProps {
    onOpen: () => void;
}

const HamburgerButton: React.FC<HamburgerButtonProps> = ({ onOpen }) => {
    const { isOpen } = useDisclosure();

    return (
        <IconButton
            onClick={onOpen}
            icon={isOpen ? <CloseIcon w={3} h={3} /> : <HamburgerIcon w={5} h={5} />}
            variant="ghost"
            aria-label="Toggle Navigation"
        />
    )
}
export default HamburgerButton;
