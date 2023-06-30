import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import SignUpForm from './SignUpForm';

jest.mock('next/router', () => require('next-router-mock'));

describe('SignUpForm', () => {
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

    it('calls the onChange handler when typing into password input', () => {
        const { getByTestId } = render(<SignUpForm />);
        const passwordInput = getByTestId('password-input') as HTMLInputElement;
        const testValue = 'password123';
        fireEvent.change(passwordInput, { target: { value: testValue } });
        expect(passwordInput.value).toBe(testValue);
    });
});

