import React, {useState} from 'react';
import {MenuItem} from '@chakra-ui/react';
import EditPostModal from "./EditPostModal";
import {PostComponent, PostView} from "../Post"

interface EditPostItemProps {
    post: PostComponent;
    updatePostHandler: React.Dispatch<React.SetStateAction<PostView>>;
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
            <EditPostModal 
                isOpen={isOpen} 
                postComponent={props.post}
                setIsOpen={setIsOpen}
                updatePostHandler={props.updatePostHandler}/>
        </>
    )
}