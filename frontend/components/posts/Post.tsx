import React from "react";
import {Avatar, Box, Button, Card, CardHeader, CardBody, CardFooter, Flex, Heading, IconButton, Menu, MenuButton, MenuItem, MenuList, Text} from "@chakra-ui/react";
import {BiComment, BiLike, BiShare} from "react-icons/bi";
import {BsThreeDotsVertical} from "react-icons/bs";
import EditPostItem from "./EditPostItem";

export interface PostView {
    User: User,
    Post: PostComponent,
    IsEditable: boolean,
}

interface User {
    Name: string,
    URL: string,
    ProfilePic: string,
}

export interface PostComponent {
    ID: number,
    CreatedAt: string,
    UpdatedAt: string,
    Content: string,
}


// TODO: Determine number of lines in text content and display Show More

export default function Post(post : PostView) {
    const notEdited = post.Post.UpdatedAt == post.Post.CreatedAt;
    const timeStamp = new Date(post.Post.UpdatedAt).toLocaleString("en-GB", {
        dateStyle: "medium",
        timeStyle: "short",
    });
    return (
        <Card bg="white" maxW="2xl" marginBlock="30px" marginInline="auto">
            <CardHeader>
                <Flex gridGap={4}>
                    <Flex 
                        flex="1" 
                        gap="4" 
                        alignItems="center" 
                        flexWrap="wrap">
                        <a href={post.User.URL}>
                            <Avatar
                            size={'md'}
                            src={post.User.ProfilePic}
                            />
                         </a>
                         <Box>
                            <a href={post.User.URL}>
                                <Heading 
                                    size="sm"
                                    _hover={{textDecoration: "underline"}}
                                    >
                                    {post.User.Name == "" ? "Anonymous User" : post.User.Name}
                                </Heading>
                            </a>
                            <Text fontSize="15px">{notEdited ? `Posted on ${timeStamp}` : `Last edited on ${timeStamp}`}</Text>
                        </Box>
                    </Flex>
                    {post.IsEditable &&
                    <Menu>
                        <MenuButton as={IconButton} variant='ghost'
                            colorScheme='gray'
                            aria-label='Options'
                            icon={<BsThreeDotsVertical />}>
                        </MenuButton>
                    <MenuList>
                        <EditPostItem post={post.Post}/>
                        <MenuItem>Delete post</MenuItem>
                    </MenuList>
                    </Menu>}
                </Flex>
            </CardHeader>
            <CardBody paddingBlock="0px">
                <Text>
                    {post.Post.Content}
                </Text>
            </CardBody>
            <CardFooter
                justify="space-between"
                flexWrap="wrap"
                paddingBlockEnd="0px"
                paddingInline="0px"
                sx={{
                    "& > button": {
                        minW: "136px"
                    }
                }}
            >
                <Button flex="1" variant="outline" leftIcon={<BiLike />}>
                    Like
                </Button>
                <Button flex="1" variant="outline" leftIcon={<BiComment />}>
                    Comment
                </Button>
                <Button flex="1" variant="outline" leftIcon={<BiShare/>}>
                    Share
                </Button>
            </CardFooter>
        </Card>
    )
}