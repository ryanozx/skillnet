import React from 'react';
import ProfileInfo from "./ProfileInfo";
import DefaultLayoutContainer from "../base/DefaultLayoutContainer";

interface ProfilePageContainerProps {
    username: string;
}

export default function ProfilePageContainer({username}: ProfilePageContainerProps) {
    return (
        <DefaultLayoutContainer>
            <ProfileInfo username={username}></ProfileInfo>
        </DefaultLayoutContainer>
    )
}