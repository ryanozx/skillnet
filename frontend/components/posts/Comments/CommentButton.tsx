import React, {useState} from "react";
import {Button} from "@chakra-ui/react";
import {BiComment} from "react-icons/bi"
import CommentsModal from "./CommentsModal";

interface CommentButtonProps {
    postID: number;
    setCommentCountHandler: React.Dispatch<React.SetStateAction<number>>;
}

export default function CommentButton(props : CommentButtonProps) {
    const [commentModalOpen, setCommentModalOpen] = useState<boolean>(false);

    const handleOpen = () => setCommentModalOpen(true);

    return (
        <>
            <Button 
                flex="1" 
                variant="outline" 
                leftIcon={<BiComment />}
                onClick={handleOpen}
            >
                Comments
            </Button>
            <CommentsModal 
                postID={props.postID}
                isOpen={commentModalOpen}
                setIsOpen={setCommentModalOpen}
                setCommentCountHandler={props.setCommentCountHandler}
            />
        </>
    )
}