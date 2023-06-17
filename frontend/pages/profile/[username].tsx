import React from "react";
import ProfileInfo from "../../components/profilePage/ProfileInfo";
import DefaultLayoutContainer from "../../components/base/DefaultLayoutContainer";
import { useRouter } from 'next/router';

interface User {
    AboutMe: string;
    Email: string;
    Name: string;
    Title: string;
    ProfilePic: string;
    Username: string;
}

export default function ProfilePage() {
    const router = useRouter();
    const { username } = router.query;
    return (
        <DefaultLayoutContainer>
            <ProfileInfo username={username as string}></ProfileInfo>
        </DefaultLayoutContainer>
    );
}