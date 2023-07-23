import React from "react";
import { useRouter } from 'next/router';
import ProfilePageContainer from "../../components/profilePage/ProfilePageContainer";



export default function ProfilePage() { 
    const router = useRouter();
    const {query} = router;
    
    return (
        <ProfilePageContainer username={query.username as string}/>
    );
}