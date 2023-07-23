import React, {useState} from "react"
import { UserMinimal } from "../../../types";
import {
    Avatar,
    Box,
    Button,
    Card,
    CardBody,
    Flex,
    Heading,
    Input,
    InputGroup,
    InputRightElement,
    Menu,
    MenuButton,
    MenuItem,
    MenuList,
    Text,
    useToast,
} from "@chakra-ui/react"
import { IconButton } from "@chakra-ui/react";
import { BsThreeDotsVertical } from "react-icons/bs";
import {BiSolidSend} from "react-icons/bi";
import DeleteCommentItem from "./DeleteCommentItem";
import axios from "axios";

export interface CommentView {
    User: UserMinimal;
    Comment: CommentComponent
    IsEditable: boolean;
    CommentCountHandler: React.Dispatch<React.SetStateAction<number>>;
}

export interface CommentComponent {
    ID: number;
    CreatedAt: string;
    UpdatedAt: string;
    Text: string;
}

export default function Comment(currComment : CommentView) {
    const [isDeleted, setIsDeleted] = useState<boolean>(false);
    const [isEditing, setIsEditing] = useState<boolean>(false);
    const [commentText, setCommentText] = useState<string>(currComment.Comment.Text); 
    const [editText, setEditText] = useState<string>(currComment.Comment.Text);
    const [editedTime, setEditedTime] = useState<string>(currComment.Comment.UpdatedAt);
    const toast = useToast();
    var notEdited = currComment.Comment.CreatedAt == editedTime;
    const baseURL = process.env.BACKEND_BASE_URL;
    const timeStamp = new Date(editedTime).toLocaleString("en-GB", {
        dateStyle: "medium",
        timeStyle: "short",
    });

    const enterEditMode = () => {
        setIsEditing(true);
        setEditText(commentText);
    }

    const exitEditMode = () => {
        setIsEditing(false);
        setEditText(commentText);
    }

    const onEditSubmit = async (e : React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        axios.patch(baseURL + "/auth/comments/" + currComment.Comment.ID, {"text": editText}, {withCredentials: true})
        .then(res => {
                setCommentText(res.data["data"]["Comment"]["Text"]);
                setEditedTime(res.data["data"]["Comment"]["UpdatedAt"]);
                exitEditMode();
        })
        .catch(err => {
            console.log(err);
            toast({
                title: "An error occurred.",
                description: err.response.data.error,
                status: "error",
                duration: 5000,
                isClosable: true,
            });
        })
    }

    return (
    <>
            {!isDeleted &&
            <Card bg="white" maxW="2xl" marginBlock="5px" marginInline="auto">
                <CardBody>
                    <Flex gridGap={4}>
                        <Flex
                            flex="1"
                            gap="4" 
                            flexWrap="wrap">
                            <a href={currComment.User.URL}>
                                <Avatar
                                size={'md'}
                                src={currComment.User.ProfilePic}
                                />
                            </a>
                            <Box>
                                <Heading 
                                    size="xs"
                                    _hover={{textDecoration: "underline"}}
                                    >
                                    <a href={currComment.User.URL}>{currComment.User.Name == "" ? "Anonymous User" : currComment.User.Name}</a>
                                </Heading>
                                <Text fontSize="13px">{notEdited ? `Posted on ${timeStamp}` : `Last edited on ${timeStamp}`}</Text>
                                {!isEditing && <Text marginBlockStart="10px">
                                    {commentText}
                                </Text>}
                                {
                                    isEditing && 
                                    <form onSubmit={e => onEditSubmit(e)} style={{flex: "1"}}>
                                        <Flex
                                            flex="1"
                                            flexWrap="wrap"
                                            gap="4"
                                        >
                                        <InputGroup>
                                            <Input 
                                                flex="1"
                                                variant="filled"
                                                placeholder="Add a comment..."
                                                value={editText}
                                                onChange={e => setEditText(e.target.value)}
                                            />
                                            <InputRightElement width="60px">
                                                <Button 
                                                    height="90%" 
                                                    variant="ghost" 
                                                    rightIcon={<BiSolidSend />}
                                                    isDisabled={editText == "" || editText == commentText}
                                                    type="submit"
                                                    />
                                            </InputRightElement>
                                        </InputGroup>
                                        <Button size="sm" onClick={exitEditMode}>Cancel</Button>
                                        </Flex>
                                    </form>
                                }
                            </Box>
                        </Flex>
                        {currComment.IsEditable &&
                        <Menu>
                            <MenuButton as={IconButton} variant='ghost'
                                colorScheme='gray'
                                aria-label='Options'
                                icon={<BsThreeDotsVertical />}>
                            </MenuButton>
                        <MenuList>
                            <MenuItem onClick={enterEditMode}>Edit comment</MenuItem>
                            <DeleteCommentItem comment={currComment.Comment} deleteHandler={setIsDeleted} commentCountHandler={currComment.CommentCountHandler}/>
                        </MenuList>
                        </Menu>}
                    </Flex>
                </CardBody>
            </Card>}
            {isDeleted &&
                <Card bg="white" maxW="2xl" marginBlock="5px" marginInline="auto">
                    <CardBody>
                        <Text align="center">This comment has been deleted.</Text>
                    </CardBody>
                </Card>
            }
        </>
    )
}