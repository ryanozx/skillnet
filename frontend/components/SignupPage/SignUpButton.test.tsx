import React from 'react';
import { render, fireEvent, act } from '@testing-library/react';
import { SignUpButton  } from './SignUpButton';
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

describe('SignUpButton', () => {
    const mockForm = {
        username: 'username',
        email: 'test@test.com',
        password: 'password',
    };

    const url = process.env.BACKEND_BASE_URL + '/signup';

    beforeEach(() => {
        mockAxios.resetHistory();
        mockRouter.mockClear();
        mockToast.mockClear();
    });

    afterEach(() => {
        mockAxios.reset();
    });

    it('should render a submit button', () => {
        const { getByRole } = render(<SignUpButton form={mockForm} />);
        const button = getByRole('button', { name: 'Sign up' });
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
        const { getByRole } = render(<SignUpButton  form={mockForm} />);
        const button = getByRole('button', { name: 'Sign up' });
        await act(async () => {
            fireEvent.click(button);
        });

        // Assertions
        expect(mockAxios.history.post[0].url).toEqual(url);
        expect(mockAxios.history.post[0].withCredentials).toEqual(true);
        expect(push).toHaveBeenCalledWith('/create-profile');
        expect(toast).toHaveBeenCalledWith({
            title: "Form submission successful.",
            description: "Account successfully created.",
            status: "success",
            duration: 5000,
            isClosable: true,
        });
    });

    it('should show error toast on request failure', async () => {
        // Setup
        const errorResponse = { message: 'Some error' };
        mockAxios.onPost(url).reply(500, errorResponse);
        const toast = jest.fn();
        mockToast.mockReturnValue(toast);

        // Render and click
        const { getByRole } = render(<SignUpButton  form={mockForm} />);
        const button = getByRole('button', { name: 'Sign up' });
        await act(async () => {
            fireEvent.click(button);
        });

        // Assertions
        expect(mockAxios.history.post[0].url).toEqual(url);
        expect(mockAxios.history.post[0].withCredentials).toEqual(true);
        expect(toast).toHaveBeenCalledWith({
            title: "An error occurred.",
            description: errorResponse.message,
            status: "error",
            duration: 5000,
            isClosable: true,
        });
    });
});


