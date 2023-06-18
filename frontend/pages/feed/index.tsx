import React from "react";
import DefaultLayoutContainer from "../../components/base/DefaultLayoutContainer";
import { useRouter } from 'next/router';
import Feed from "../../components/posts/Feed";

export default function ProfilePage() {
    const router = useRouter();
    return (
        <DefaultLayoutContainer>
            <Feed></Feed>
        </DefaultLayoutContainer>
    );
}