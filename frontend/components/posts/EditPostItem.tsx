import React, {useState} from 'react';
import {MenuItem} from '@chakra-ui/react';
import EditPostModel from "./EditPostModel";
import {PostComponent} from "./Post"

interface EditPostItemProps {
    post: PostComponent
}

export default function EditPostItem(props: EditPostItemProps) {
    const [isOpen, setIsOpen] = useState<boolean>(false)

    const handleOpen = () => setIsOpen(true)

    return (
        <>
            <MenuItem
                onClick={handleOpen}
            >Edit post
            </MenuItem>
            <EditPostModel 
                isOpen={isOpen} 
                postComponent={props.post}
                setIsOpen={setIsOpen}/>
        </>
    )
}