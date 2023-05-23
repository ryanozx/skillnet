import React from "react";
// import ProfileInfo from "../../components/ProfilePage/ProfileInfo";
import ProfileInfo from "../../components/profilePage/ProfileInfo";
import LayoutContainer from "../../components/base/LayoutContainer";
import { useRouter } from 'next/router';

export default function ProfilePage() {
    const router = useRouter();
    const { user_id } = router.query;
    return (
        <LayoutContainer>
            <ProfileInfo user_id={user_id}></ProfileInfo>
        </LayoutContainer>
    );

}