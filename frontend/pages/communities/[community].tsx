import React from "react";
import { useRouter } from 'next/router';
import CommunityPageContainer from "../../components/communityPage/CommunityPageContainer";
import { escapeHtml } from "../../types";

export default function ProfilePage() { 
    const router = useRouter();
    const {query} = router;
    
    return (
        <CommunityPageContainer communityName={escapeHtml(query.community as string)}/>
    );
}