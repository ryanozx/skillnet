import React from 'react';
import BasicInfo from './BasicInfo';
import AboutMe from './AboutMe';
import { User } from '../../types';

interface InfoSectionProps {
    user: User;
    username: string;
    setUser?: React.Dispatch<React.SetStateAction<User>>;
}

export default function InfoSection(props : InfoSectionProps) {
    return (
        <>
            <BasicInfo 
                user={props.user}
                username={props.username}
                setUser={props.setUser}/>
            <AboutMe aboutMe={props.user.AboutMe}></AboutMe>
        </>
    )
}