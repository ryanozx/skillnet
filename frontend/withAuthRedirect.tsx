import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import axios from 'axios';
import LoadingScreen from './components/base/LoadingScreen';

const loadingText = "Please be patient while we load our demo, we are using a free server ...";

export function preventAuthAccess(Component: any) {
    
    return function AuthenticatedComponent(props: any) {
        const router = useRouter();
        const [loading, setLoading] = useState(true);
        const base_url = process.env.BACKEND_BASE_URL;
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
            return <LoadingScreen loadingText={loadingText}/>
        }

        return <Component {...props} />;
    }
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
            // return <LoadingScreen loadingText={loadingText}/>
            return null;
        }

        return <Component {...props} />;
    };
}

