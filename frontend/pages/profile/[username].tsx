import React from "react";
import {useEffect} from "react";
import ProfileInfo from "../../components/profilePage/ProfileInfo";
import DefaultLayoutContainer from "../../components/base/DefaultLayoutContainer";
import { useRouter } from 'next/router';

export default function ProfilePage() {
    const router = useRouter();
    const {isReady, query} = router;
    
    console.log(query.username)
    
    return (
        <DefaultLayoutContainer>
            <ProfileInfo ownProfile={false}></ProfileInfo>
        </DefaultLayoutContainer>
    );
}