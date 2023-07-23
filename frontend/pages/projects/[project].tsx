import {useRouter} from "next/router";
import ProjectPageContainer from "../../components/projectPage/ProjectPageContainer";

export default function ProjectPage() {
    const router = useRouter()
    const projectID = parseInt(router.query.project as string, 10)

    return (
        <ProjectPageContainer projectID={projectID} />
    )
}