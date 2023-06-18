import React from "react";
import ProfileInfo from "../../../components/profilePage/ProfileInfo";
import DefaultLayoutContainer from "../../../components/base/DefaultLayoutContainer";

export default function ProfilePage() {

    // must do something here to get the username from sessionId token

    return (
        <DefaultLayoutContainer>
            <ProfileInfo ownProfile={true}></ProfileInfo>
        </DefaultLayoutContainer>
    );
}