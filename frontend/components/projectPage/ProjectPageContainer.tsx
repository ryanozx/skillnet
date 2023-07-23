import React, {useState} from 'react';
import DefaultLayoutContainer from "../base/DefaultLayoutContainer";
import Feed from '../posts/Feed';
import ProjectInfo from './ProjectInfo';

interface ProjectPageContainerProps {
    projectID: number
}

export default function ProjectPageContainer(props: ProjectPageContainerProps) {
    const [projectLoaded, setProjectLoaded] = useState<boolean>(false);
    const [communityID, setCommunityID] = useState<number>(0);
    return (
        <DefaultLayoutContainer>
            {props.projectID && <ProjectInfo ProjectID={props.projectID} setProjectLoaded={setProjectLoaded} setCommunityID={setCommunityID}/>}
            {projectLoaded && <Feed AllowPostAdd={true} ProjectID={props.projectID} CommunityID={communityID}/>}
        </DefaultLayoutContainer>
    )
}