import React from 'react';
import { screen, render, fireEvent } from '@testing-library/react';
import SignUpForm from './SignUpForm';

jest.mock('next/router', () => require('next-router-mock'));

describe('SignUpForm', () => {

    it('should render EmailInput', () => {
        const { getByTestId } = render(<SignUpForm />);
        const emailInput = getByTestId('email-input');
        expect(emailInput).toBeInTheDocument();
    });

    it('should render UsernameInput', () => {
        const { getByTestId } = render(<SignUpForm />);
        const usernameInput = getByTestId('username-input');
        expect(usernameInput).toBeInTheDocument();
    });

    it('should render PasswordInput', () => {
        const { getByTestId } = render(<SignUpForm />);
        const passwordInput = getByTestId('password-input');
        expect(passwordInput).toBeInTheDocument();
    });

    it('should render SignUpButton', () => {
        const { getByTestId } = render(<SignUpForm />);
        const signUpButton = getByTestId('signup-button');
        expect(signUpButton).toBeInTheDocument();
    });

    it('should render LoginRedirect', () => {
        const { getByTestId } = render(<SignUpForm />);
        const loginRedirectText = getByTestId('redirect-text');
        const loginReddirectLink = getByTestId('login-link');
        expect(loginRedirectText).toBeInTheDocument();
        expect(loginReddirectLink).toBeInTheDocument();
    });

    it('updates email input value when typing', () => {
        const { getByTestId } = render(<SignUpForm />);
        const emailInput = getByTestId('email-input') as HTMLInputElement;
        fireEvent.change(emailInput, { target: { value: 'test@email.com' } });
        expect(emailInput.value).toBe('test@email.com');
    });

    it('updates username input value when typing', () => {
        const { getByTestId } = render(<SignUpForm />);
        const usernameInput = getByTestId('username-input') as HTMLInputElement;
        fireEvent.change(usernameInput, { target: { value: 'username123' } });
        expect(usernameInput.value).toBe('username123');
    });

    it('updates password input value when typing', () => {
        const { getByTestId } = render(<SignUpForm />);
        const passwordInput = getByTestId('password-input') as HTMLInputElement;
        const testValue = 'password123';
        fireEvent.change(passwordInput, { target: { value: testValue } });
        expect(passwordInput.value).toBe(testValue);
    });
});

