import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import axios from 'axios';

export function preventAuthAccess(Component: any) {
    
    return function AuthenticatedComponent(props: any) {
        const router = useRouter();
        const [loading, setLoading] = useState(true);
        const base_url = process.env.NEXT_PUBLIC_BACKEND_BASE_URL
        useEffect(() => {
            axios.get(base_url + '/auth/user', { withCredentials: true })
                .then((res) => {
                    // if we get a successful response, the user is logged in, so redirect
                    const {Username} = res.data.data;
                    router.push(`/feed`);
                })
                .catch((error) => {
                    // if there was an error, we couldn't get user info, which means the user is not logged in
                    setLoading(false);
                });
        }, []);

        if (loading) {
            return null; // you can return a loading spinner here if you want
        }

        return <Component {...props} />;
    };
}

export function requireAuth(Component: any) {
    return function AuthenticatedComponent(props: any) {
        const router = useRouter();
        const [loading, setLoading] = useState(true);
        const base_url = process.env.BACKEND_BASE_URL;
        useEffect(() => {
            axios.get(base_url + '/auth/user', { withCredentials: true })
                .then((res) => {
                    // if we get a successful response, the user is logged in, so redirect
                    setLoading(false);
                })
                .catch((error) => {
                    // if there was an error, we couldn't get user info, which means the user is not logged in
                    
                    router.push(`/`);
                });
        }, []);

        if (loading) {
            return null; // you can return a loading spinner here if you want
        }

        return <Component {...props} />;
    };
}

