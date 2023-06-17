import React from 'react';
import ProfileInfo from "./ProfileInfo";
import DefaultLayoutContainer from "../base/DefaultLayoutContainer";

export default function ProfilePageContainer() {
    return (
        <DefaultLayoutContainer>
            <ProfileInfo username="ivyy-poison"></ProfileInfo>
        </DefaultLayoutContainer>
    )
}