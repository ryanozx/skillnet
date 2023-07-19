import React, {useState, useEffect} from "react";
import {
    Avatar,
    Button,
    Card,
    CardBody,
    Flex,
    Input,
    InputGroup,
    InputRightElement,
    Link,
    useToast,
} from "@chakra-ui/react"
import {
    BiSolidSend
} from "react-icons/bi"
import axios from "axios";
import { CommentView } from "./Comment";

interface CreateCommentProps {
    postID: number;
    addCommentHandler: (comment : CommentView) => void;
    setCommentCountHandler: React.Dispatch<React.SetStateAction<number>>;
}

export default function CreateCommentCard(props : CreateCommentProps) {
    const [needUpdate, setNeedUpdate] = useState<boolean>(true);
    const [profilePic, setProfilePic ] = useState<string>("");
    const [username, setUsername ] = useState<string>("");
    const [text, setText] = useState<string>("");
    const toast = useToast();

    const baseURL = process.env.BACKEND_BASE_URL;

    useEffect(() => {
        console.log("Initial retrieval")
        console.log("Retrieving profile...")
        axios.get(baseURL + '/auth/user', { withCredentials: true })
            .then((res) => {
                const { ProfilePic, Username } = res.data.data;
                setProfilePic(ProfilePic);
                setUsername(Username);
                setNeedUpdate(false); // reset the flag after the data is updated
                console.log("profile retrieved")
            })
            .catch((err) => {
                console.log(err);
        });
    }, []);

    const onSubmit = async (e : React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        axios.post(baseURL + "/auth/comments?post=" + props.postID, {"text": text}, {withCredentials: true})
        .then(res => {
                props.addCommentHandler(res.data["data"]["Comment"])
                props.setCommentCountHandler(res.data["data"]["CommentCount"])
                console.log("Successfully created comment")
                setText("");
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
        <Card bg="transparent" variant="unstyled" width={"100%"} marginInline="auto" marginBlockEnd="10px">
            <CardBody>
                <Flex
                    gap="4" 
                    alignItems="center"
                    flexWrap="wrap"
                    >
                    <Link href={`/profile/${username}`}>
                        <Avatar
                            size="sm"
                            src={profilePic}
                        />
                    </Link>
                    <form onSubmit={e => onSubmit(e)} style={{flex: "1"}}>
                        <InputGroup>
                            <Input 
                                flex="1"
                                variant="filled"
                                placeholder="Add a comment..."
                                value={text}
                                onChange={e => setText(e.target.value)}
                            />
                            <InputRightElement width="60px">
                                <Button 
                                    height="90%" 
                                    variant="ghost" 
                                    rightIcon={<BiSolidSend />}
                                    isDisabled={text == ""}
                                    type="submit"
                                    />
                            </InputRightElement>
                        </InputGroup>
                    </form>
                </Flex>
            </CardBody>
    </Card>
    )
}