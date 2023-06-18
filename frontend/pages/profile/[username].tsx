import React from "react";
import ProfileInfo from "../../components/profilePage/ProfileInfo";
import DefaultLayoutContainer from "../../components/base/DefaultLayoutContainer";
import { useRouter } from 'next/router';
import ProfilePageContainer from "../../components/profilePage/ProfilePageContainer";



export default function ProfilePage() { 
    const router = useRouter();
    const { username } = router.query;
    return (
        <ProfilePageContainer username={username as string}/>
    );
}