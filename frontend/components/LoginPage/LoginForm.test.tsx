import React from 'react';
import { screen, render, fireEvent } from '@testing-library/react';
import LoginForm from './LoginForm';

jest.mock('next/router', () => require('next-router-mock'));

describe('LoginForm', () => {
    it('should render UsernameInput', () => {
        const { getByTestId } = render(<LoginForm />);
        const usernameInput = getByTestId('username-input');
        expect(usernameInput).toBeInTheDocument();
        screen.debug()
    });

    it('should render PasswordInput', () => {
        const { getByTestId } = render(<LoginForm />);
        const passwordInput = getByTestId('password-input');
        expect(passwordInput).toBeInTheDocument();
    });

    it('should render LoginButton', () => {
        const { getByTestId } = render(<LoginForm />);
        const signUpButton = getByTestId('login-button');
        expect(signUpButton).toBeInTheDocument();
    });

    it('should render SignupRedirect', () => {
        const { getByTestId } = render(<LoginForm />);
        const loginRedirectText = getByTestId('redirect-text');
        const loginReddirectLink = getByTestId('signup-link');
        expect(loginRedirectText).toBeInTheDocument();
        expect(loginReddirectLink).toBeInTheDocument();
    });

    it('updates username input value when typing', () => {
        const { getByTestId } = render(<LoginForm />);
        const usernameInput = getByTestId('username-input') as HTMLInputElement;
        fireEvent.change(usernameInput, { target: { value: 'username123' } });
        expect(usernameInput.value).toBe('username123');
    });

    it('updates password input value when typing', () => {
        const { getByTestId } = render(<LoginForm />);
        const passwordInput = getByTestId('password-input') as HTMLInputElement;
        const testValue = 'password123';
        fireEvent.change(passwordInput, { target: { value: testValue } });
        expect(passwordInput.value).toBe(testValue);
    });
});
