import React, {useState} from 'react';
import DefaultLayoutContainer from "../base/DefaultLayoutContainer";
import Feed from '../posts/Feed';
import CommunityInfo from './CommunityInfo';

interface CommunityPageContainerProps {
    communityName: string
}

export default function CommunityPageContainer(props: CommunityPageContainerProps) {
    const [communityLoaded, setCommunityLoaded] = useState<boolean>(false);
    const [communityID, setCommunityID] = useState<number>(0);
    return (
        <DefaultLayoutContainer>
            <CommunityInfo communityName={props.communityName} setCommunityID={setCommunityID} setCommunityLoaded={setCommunityLoaded}/>
            {communityLoaded && <Feed AllowPostAdd={false} CommunityID={communityID} isGlobal={false}/>}
        </DefaultLayoutContainer>
    )
}