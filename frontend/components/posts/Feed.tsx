import React, { useState, useEffect, useRef } from 'react';
import axios from "axios";
import {Box, Text} from "@chakra-ui/react";
import Post, {PostView} from "./Post";
import InfiniteScroll from "react-infinite-scroll-component";
import CreatePostCard from "./CreatePostCard";

interface FeedProps {
    AllowPostAdd: boolean
    CommunityID?: number
    ProjectID?: number
    isGlobal: boolean
}

export default function Feed(props: FeedProps) {
    const [posts, setPosts] = useState<React.JSX.Element[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState(null);
    const base_url = process.env.BACKEND_BASE_URL;
    const [url, setURL] = useState(base_url + "/auth/posts")
    const [noMorePosts, setNoMorePosts] = useState<boolean>(false);
    const [initialUrlUpdated, setInitialUrlUpdated] = useState<boolean>(false);

    useEffect(() => {
        if (props.ProjectID && props.ProjectID != 0) {
            setURL(url + "?project=" + props.ProjectID)
            setInitialUrlUpdated(true);
        }
        else if (props.CommunityID && props.CommunityID != 0) {
            setURL(url + "?community=" + props.CommunityID)
            setInitialUrlUpdated(true);
        }
        else if (props.isGlobal) {
            setInitialUrlUpdated(true);
        }
    }, [props.ProjectID, props.CommunityID, props.isGlobal])

    useEffect(() => {
        if (initialUrlUpdated)
        {
            updateFeed();
        }
    }, [initialUrlUpdated])

    const updateFeed = async() => {
        if (!isLoading) {
            setIsLoading(true);
            setError(null);

            console.log(url)
            
            const fetchData = axios.get(url, {withCredentials: true});
            fetchData
            .then((response) => {
                if (response.data["data"]["Posts"] != null)
                {
                    if (response.data["data"]["Posts"].length == 0) {
                        setNoMorePosts(true);
                    }
                    setPosts([...posts, ...response.data["data"]["Posts"].map((postdata : PostView) => <Post key={postdata.Post.ID} {...postdata}/>)]);
                } else {
                    setNoMorePosts(true);
                }
                setURL(response.data["data"]["NextPageURL"]);
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
        setPosts([<Post key={post.Post.ID} {...post}/>, ...posts]);
    }

    return (<Box>
        {props.AllowPostAdd && props.CommunityID && props.ProjectID && <CreatePostCard addPostHandler={addPost} communityID={props.CommunityID} projectID={props.ProjectID}/>}
        <InfiniteScroll
            dataLength={posts.length}
            next={updateFeed}
            hasMore={!noMorePosts}
            loader={
                <Box paddingBlock="10px">
                    <Text textAlign="center">Loading...</Text>
                </Box>
            }
            endMessage={
                <Box paddingBlock="10px">
                    <Text textAlign="center">No more posts to load.</Text>
                </Box>}
            >
            {posts}
        </InfiniteScroll>

        {error && 
        <Box paddingBlock="10px">
            <Text textAlign="center">There was an error in loading the posts.</Text>
        </Box>}
    </Box>)
}
