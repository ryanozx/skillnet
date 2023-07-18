import {
    Drawer,
    DrawerBody,
    DrawerContent,
    DrawerCloseButton,
    DrawerOverlay,
    DrawerHeader,
    Link,
} from '@chakra-ui/react';
import React from 'react';
import SideBar from '../SideBar/SideBar';

interface DrawerMenuProps {
    isOpen: boolean;
    onClose: () => void;
    btnRef: React.RefObject<any>;
}

const DrawerMenu: React.FC<DrawerMenuProps> = ({ isOpen, onClose, btnRef }) => {
    return (
        <Drawer
            isOpen={isOpen}
            placement='left'
            onClose={onClose}
            finalFocusRef={btnRef}
        >
            <DrawerOverlay />
            <DrawerContent>
                <DrawerCloseButton />
                <Link href="/feed">
                    <DrawerHeader>SkillNet</DrawerHeader>
                </Link>
                <DrawerBody>
                    <SideBar/>
                </DrawerBody>
            </DrawerContent>
        </Drawer>
    )
}
export default DrawerMenu;
