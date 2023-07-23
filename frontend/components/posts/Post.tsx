import React, {useState} from "react";
import {Avatar, Box, Button, Card, CardHeader, CardBody, CardFooter, Flex, Heading, IconButton, Menu, MenuButton, MenuItem, MenuList, Text} from "@chakra-ui/react";
import {BiShare} from "react-icons/bi";
import {BsThreeDotsVertical} from "react-icons/bs";
import EditPostItem from "./EditPost/EditPostItem";
import DeletePostItem from "./DeletePostItem";
import LikeButton from "./LikeButton";
import CommentButton from "./Comments/CommentButton";
import {UserMinimal} from "./../../types";
import { ProfileButtonProps } from "../base/NavBar/ProfileButton";

export interface PostAndComments{
    post: PostView,
    profile: ProfileButtonProps
}

export interface PostView {
    User: UserMinimal,
    Post: PostComponent,
    IsEditable: boolean,
    Liked: boolean,
    LikeCount: number,
    CommentCount: number,
}

export interface PostComponent {
    ID: number,
    CreatedAt: string,
    UpdatedAt: string,
    Content: string,
}

// TODO: Determine number of lines in text content and display Show More

export default function Post(post : PostView) {
    const [currPost, setCurrPost] = useState<PostView>(post);
    const [isDeleted, setIsDeleted] = useState<boolean>(false);
    const [isLiked, setIsLiked] = useState<boolean>(post.Liked);
    const [likeCount, setLikeCount] = useState<number>(post.LikeCount);
    const [commentCount, setCommentCount] = useState<number>(post.CommentCount)
    const notEdited = currPost.Post.UpdatedAt == currPost.Post.CreatedAt;
    const timeStamp = new Date(currPost.Post.UpdatedAt).toLocaleString("en-GB", {
        dateStyle: "medium",
        timeStyle: "short",
    });
    return (
        <>
            {!isDeleted &&
            <Card bg="white" maxW="2xl" marginBlock="30px" marginInline="auto">
                <CardHeader>
                    <Flex gridGap={4}>
                        <Flex 
                            flex="1" 
                            gap="4" 
                            alignItems="center" 
                            flexWrap="wrap">
                            <a href={currPost.User.URL}>
                                <Avatar
                                size={'md'}
                                src={currPost.User.ProfilePic}
                                />
                            </a>
                            <Box>
                                <Heading 
                                    size="sm"
                                    _hover={{textDecoration: "underline"}}
                                    >
                                    <a href={currPost.User.URL}>{(currPost.User.Name == "" || !currPost.User.Name) ? "Anonymous User" : currPost.User.Name}</a>
                                </Heading>
                                <Text fontSize="15px">{notEdited ? `Posted on ${timeStamp}` : `Last edited on ${timeStamp}`}</Text>
                            </Box>
                        </Flex>
                        {currPost.IsEditable &&
                        <Menu>
                            <MenuButton as={IconButton} variant='ghost'
                                colorScheme='gray'
                                aria-label='Options'
                                icon={<BsThreeDotsVertical />}>
                            </MenuButton>
                        <MenuList>
                            <EditPostItem post={currPost.Post} updatePostHandler={setCurrPost}/>
                            <DeletePostItem post={currPost.Post} deleteHandler={setIsDeleted}/>
                        </MenuList>
                        </Menu>}
                    </Flex>
                </CardHeader>
                <CardBody paddingBlock="0px">
                    <Text>
                        {currPost.Post.Content}
                    </Text>
                </CardBody>
                <CardFooter
                    justify="space-between"
                    flexWrap="wrap"
                    paddingBlockEnd="0px"
                    paddingInline="20px"
                >
                    <Text color="gray">{likeCount}{likeCount === 1 ? " like" : " likes"}</Text>
                    <Text color="gray">{commentCount}{commentCount === 1 ? " comment" : " comments"}</Text>
                </CardFooter>
                <CardFooter
                    justify="space-between"
                    flexWrap="wrap"
                    paddingBlock="5px"
                    paddingInline="0px"
                    sx={{
                        "& > button": {
                            minW: "136px"
                        }
                    }}
                >
                    <LikeButton 
                        Liked={isLiked}
                        PostID={currPost.Post.ID}  
                        SetLikeCountHandler={setLikeCount}
                        SetLikedHandler={setIsLiked} 
                        />
                    <CommentButton 
                        postID={currPost.Post.ID}
                        setCommentCountHandler={setCommentCount}
                    />
                    <Button flex="1" variant="outline" leftIcon={<BiShare/>}>
                        Share
                    </Button>
                </CardFooter>
            </Card>}
            {isDeleted &&
                <Card bg="white" maxW="2xl" marginBlock="30px" marginInline="auto">
                    <CardBody>
                        <Text align="center">This post has been deleted.</Text>
                    </CardBody>
                </Card>
            }
        </>
    )
}