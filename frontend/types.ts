// types.ts
export type User = {
    // Define your user properties and types here
    // Example:
    id: number;
    name: string;
} | null;
  
export interface UserState {
    loading: boolean;
    isLoggedIn: boolean;
    user: User;
    error: string;
}
  