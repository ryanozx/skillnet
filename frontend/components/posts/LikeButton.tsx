import React from "react";
import axios from "axios";
import {Button, useToast} from "@chakra-ui/react";
import {AiOutlineLike, AiFillLike} from "react-icons/ai";

interface LikeProps {
    PostID: number,
    Liked: boolean,
    SetLikedHandler: React.Dispatch<React.SetStateAction<boolean>>,
    SetLikeCountHandler: React.Dispatch<React.SetStateAction<number>>,
}

export default function LikeButton(props : LikeProps) {
    const baseURL = process.env.BACKEND_BASE_URL;
    const likeURL = baseURL + "/auth/likes/" + props.PostID.toString();
    const toast = useToast();

    const postLike = () => {
        axios.post(likeURL, {}, {withCredentials: true})
        .then(res => {
            props.SetLikeCountHandler(res.data["data"]["LikeCount"])
            props.SetLikedHandler(true);
            console.log("Liked post %d", props.PostID)
        })
        .catch(err => {
            console.log(err);
            toast({
                title: "Failed to like post",
                description: err.response.data.error,
                status: "error",
                duration: 5000,
                isClosable: true,
            });
        })
    }

    const deleteLike = () => {
        axios.delete(likeURL, {withCredentials: true})
        .then(res => {
            props.SetLikeCountHandler(res.data["data"]["LikeCount"])
            props.SetLikedHandler(false);
            console.log("Unliked post %d", props.PostID)
        })
        .catch(err => {
            console.log(err);
            toast({
                title: "Failed to unlike post.",
                description: err.response.data.error,
                status: "error",
                duration: 5000,
                isClosable: true,
            });
        })
    }

    return (
        <Button 
            flex="1" 
            leftIcon={
                props.Liked 
                ? <AiFillLike color="blue"/> 
                : <AiOutlineLike />}
            onClick={props.Liked ? deleteLike : postLike}
            variant="outline"
            >
            {props.Liked ? "Liked" : "Like"}
        </Button>
    )
}
