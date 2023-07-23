import React from 'react';
import BasicInfo from './BasicInfo';
import AboutMe from './AboutMe';
import { User } from '../../types';

interface InfoSectionProps {
    user: User;
    setUser?: React.Dispatch<React.SetStateAction<User>>;
}



export default function InfoSection({user, setUser}: InfoSectionProps) {
    return (
        <>
            <BasicInfo 
                name={user.Name}
                username={user.Username}
                title={user.Title}
                profilePic={user.ProfilePic}
                aboutMe={user.AboutMe}
                setUser={setUser}/>
            <AboutMe aboutMe={user.AboutMe}></AboutMe>
        </>
    )
}