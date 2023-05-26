import React from "react";
import ProfileInfo from "../../components/profilePage/ProfileInfo";
import LayoutContainer from "../../components/base/LayoutContainer";
import { useRouter } from 'next/router';

export default function ProfilePage() {
    const router = useRouter();
    const { username } = router.query;
    return (
        <LayoutContainer>
            <ProfileInfo username={username}></ProfileInfo>
        </LayoutContainer>
    );

}