import React from "react";
import { useRouter } from 'next/router';
import ProfilePageContainer from "../../components/profilePage/ProfilePageContainer";
import { escapeHtml } from "../../types";



export default function ProfilePage() { 
    const router = useRouter();
    const {query} = router;
    
    return (
        <ProfilePageContainer username={escapeHtml(query.username as string)}/>
    );
}