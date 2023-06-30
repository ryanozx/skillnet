import React from 'react';
import { render, fireEvent, act } from '@testing-library/react';
import { LoginButton } from './LoginButton';
import axios from 'axios';
import { useRouter } from 'next/router';
import MockAdapter from "axios-mock-adapter";
import { useToast } from '@chakra-ui/react';

jest.mock('@chakra-ui/react', () => ({
    ...jest.requireActual('@chakra-ui/react'), // keep all other real implementations
    useToast: jest.fn().mockImplementation(() => jest.fn()), // replace useToast with a mock function
}));

jest.mock('next/router', () => ({
    useRouter: jest.fn() // mock useRouter
}));

const mockAxios = new MockAdapter(axios);
const mockRouter = useRouter as jest.Mock;
const mockToast = useToast as jest.Mock;

const url = process.env.BACKEND_BASE_URL + '/login';
const username = 'username';
const password = 'password';

describe('LoginButton', () => {

    beforeEach(() => {
        mockAxios.resetHistory(); // Clears all history of requests
        mockRouter.mockClear(); // Clears all instances of this mock function
        mockToast.mockClear(); // Clears all instances of this mock function
    });
    
    afterEach(() => {
        mockAxios.reset(); // Removes any mock handlers
    });
    

    it('should render a login button', () => {
        const { getByRole } = render(<LoginButton username={username} password={password} />);
        const button = getByRole('button', { name: 'Sign in' });
        expect(button).toBeInTheDocument();
    });

    it('should make post request, navigate and show success toast on click', async () => {
        // Setup
        mockAxios.onPost(url).reply(200);
        const push = jest.fn();
        mockRouter.mockReturnValue({ push });
        const toast = jest.fn();
        mockToast.mockReturnValue(toast);
      
        // Render and click
        const { getByRole } = render(<LoginButton username={username} password={password} />);
        const button = getByRole('button', { name: 'Sign in' });
        await act(async () => {
            fireEvent.click(button);
        });

        // Assertions
        expect(mockAxios.history.post[0].url).toEqual(url);
        expect(mockAxios.history.post[0].withCredentials).toEqual(true);
        expect(push).toHaveBeenCalledWith('/feed');
        expect(toast).toHaveBeenCalledWith({
            title: "Form submission successful.",
            description: "We've received your form data.",
            status: "success",
            duration: 5000,
            isClosable: true,
        });
    });

    it('should show error toast on request failure', async () => {
        // Setup
        const errorResponse = { error: 'Some error' };
        mockAxios.onPost(url).reply(500, errorResponse);
        const toast = jest.fn();
        mockToast.mockReturnValue(toast);
      
        // Render and click
        const { getByRole } = render(<LoginButton username={username} password={password} />);
        const button = getByRole('button', { name: 'Sign in' });
        await act(async () => {
            fireEvent.click(button);
        });

        // Assertions
        expect(mockAxios.history.post[0].url).toEqual(url);
        expect(mockAxios.history.post[0].withCredentials).toEqual(true);
        expect(toast).toHaveBeenCalledWith({
            title: "An error occurred.",
            description: errorResponse.error,
            status: "error",
            duration: 5000,
            isClosable: true,
        });
    });
});
