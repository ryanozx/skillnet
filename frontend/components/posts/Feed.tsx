import React, { useState, useEffect, useRef } from 'react';
import axios from "axios";
import {Box, Text} from "@chakra-ui/react";
import Post, {PostView} from "./Post";
import InfiniteScroll from "react-infinite-scroll-component";
import CreatePostCard from "./CreatePostCard";

export default function Feed() {
    const [posts, setPosts] = useState<React.JSX.Element[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState(null);
    const [url, setURL] = useState("http://localhost:8080/auth/posts")

    const updateFeed = async() => {
        if (!isLoading) {
            console.log("Fetching data");
            setIsLoading(true);
            setError(null);

            const fetchData = axios.get(url, {withCredentials: true});
            fetchData
            .then((response) => {
                setPosts([...posts, ...response.data["data"]["Posts"].map((postdata : PostView) => <Post key={postdata.Post.ID} {...postdata}/>)])
                setURL(response.data["data"]["NextPageURL"])
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

    const addPost = (post : PostView) => {
        setPosts([<Post key={post.Post.ID} {...post} />, ...posts]);
    }

    useEffect(() => {
        updateFeed();
    }, []);

    return (<Box>
        <CreatePostCard addPostHandler={addPost}/>
        <InfiniteScroll
            dataLength={posts.length}
            next={updateFeed}
            hasMore={url != "http://localhost:8080/auth/posts?cutoff=1"}
            loader={<Box paddingBlock="10px">
            <Text textAlign="center">Loading...</Text>
        </Box>}
            endMessage={
                <Box paddingBlock="10px">
                    <Text textAlign="center">No more posts to load.</Text>
                </Box>}>
            {posts}
        </InfiniteScroll>
        {error && 
        <Box paddingBlock="10px">
            <Text textAlign="center">There was an error in loading the posts.</Text>
        </Box>}
    </Box>)
}
