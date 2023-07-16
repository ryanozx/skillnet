import React, {useState, useEffect} from 'react';
import {
    Box,
    Divider,
    Modal,
    ModalContent,
    ModalOverlay,
    Text,
} from "@chakra-ui/react";
import InfiniteScroll from 'react-infinite-scroll-component';
import CreateCommentCard from './CreateCommentCard';
import {CommentView} from "./Comment";
import Comment from "./Comment";
import axios from "axios";

interface CommentModelProps {
    isOpen: boolean;
    setIsOpen: React.Dispatch<React.SetStateAction<boolean>>;
    postID: number;
    setCommentCountHandler: React.Dispatch<React.SetStateAction<number>>;
}

export default function CommentsModel(props: CommentModelProps) {
    const [comments, setComments] = useState<React.JSX.Element[]>([]);
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [error, setError] = useState(null);
    const baseURL = process.env.BACKEND_BASE_URL;
    const [url, setURL] = useState<string>(baseURL + "/auth/comments?post=" + props.postID)

    const onClose = () => {
        props.setIsOpen(false);
        setURL(baseURL + "/auth/comments?post=" + props.postID);
        setComments([]);
    }

    const updateCommentsFeed = async() => {
        if (!isLoading) {
            console.log("Fetching data");
            setIsLoading(true);
            setError(null);

            const fetchData = axios.get(url, {withCredentials: true});
            fetchData
            .then((response) => {
                if (response.data["data"]["Comments"] != null)
                {
                    setComments([...comments, ...response.data["data"]["Comments"].map((commentData : CommentView) => <Comment key={commentData.Comment.ID} {...commentData} CommentCountHandler={props.setCommentCountHandler}/>)]);
                }
                console.log(response.data["data"]["NextPageURL"]);
                setURL(response.data["data"]["NextPageURL"]);
                console.log(url);
            })
            .catch((error) => {
                console.log(error);
                setError(error);
            })
            .finally(() =>
                setIsLoading(false)
            )
        }
    }

    const addComment = (comment : CommentView) => {
        setComments([ <Comment key={comment.Comment.ID} {...comment} CommentCountHandler={props.setCommentCountHandler}/>, ...comments])
    }

    useEffect(() => {
        props.isOpen && updateCommentsFeed();
    }, [props.isOpen]);

    return (<Modal 
        isOpen={props.isOpen}
        onClose={onClose}
        >
            <ModalOverlay />
            <ModalContent padding="20px" maxW="xl">
                <CreateCommentCard 
                    postID={props.postID} 
                    addCommentHandler={addComment}
                    setCommentCountHandler={props.setCommentCountHandler}
                    />
                <Divider borderColor="gray.400" />
                <InfiniteScroll
                    height={window.innerHeight * 0.6}
                    dataLength={comments.length}
                    next={updateCommentsFeed}
                    hasMore={url != baseURL + "/auth/comments?post=" + props.postID + "&cutoff=0"}
                    loader={
                        <Box paddingBlock="10px">
                            <Text textAlign="center">Loading...</Text>
                        </Box>
                    }
                    endMessage={
                        <Box paddingBlock="10px">
                            <Text textAlign="center">{comments.length == 0 ? "Be the first to comment!" : "No more comments to load."}</Text>
                        </Box>}>
                    {comments}
                </InfiniteScroll>
                {error && 
                <Box paddingBlock="10px">
                    <Text textAlign="center">There was an error in loading the comments.</Text>
                </Box>}
            </ModalContent>
    </Modal>)
}

