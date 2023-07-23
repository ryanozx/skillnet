// types.ts
// Contains all the types used in the frontend

export interface User {
    AboutMe: string;
    Email: string;
    Name: string;
    Title: string;
    ProfilePic: string;
    Username: string;
    ShowAboutMe: boolean;
    ShowTitle: boolean;
}

export interface Projects {
    ProjectList: ProjectMinimal[];
    NextPageURL: string;
}

export interface ProjectMinimal {
    ID: number,
    Name: string,
    Community: string,
    URL: string,
}

export interface UserMinimal {
    Name: string,
    URL: string,
    ProfilePic: string,
}