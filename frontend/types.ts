// types.ts
// Contains all the types used in the frontend

export interface User {
    AboutMe?: string;
    Email: string;
    Name?: string;
    Title?: string;
    ProfilePic?: string;
    Username: string;
    Projects: ProjectView[];
}

export interface Projects {
    ProjectList: ProjectView[];
    NextPageURL: string;
}

export interface ProjectView {
    logo: string;
    name: string;
    category: string;
}

export interface UserMinimal {
    Name: string,
    URL: string,
    ProfilePic: string,
}

export interface EditableUserInfo {
    name: string;
    username: string;
    title: string;
    profilePic: string;
    aboutMe: string;
}
