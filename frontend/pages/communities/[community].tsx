import React from "react";
import { useRouter } from 'next/router';
import CommunityPageContainer from "../../components/communityPage/CommunityPageContainer";


export default function ProfilePage() { 
    const router = useRouter();
    const {query} = router;
    
    return (
        <CommunityPageContainer communityName={query.community as string}/>
    );
}