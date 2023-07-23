import React from "react";
import DefaultLayoutContainer from "../../components/base/DefaultLayoutContainer";
import Feed from "../../components/posts/Feed";

export default function ProfilePage() {
    return (
        <DefaultLayoutContainer>
            <Feed AllowPostAdd={false}/>
        </DefaultLayoutContainer>
    );
}