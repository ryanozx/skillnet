import React from 'react';
import { screen, render, fireEvent, getByTestId, waitFor } from '@testing-library/react';
import { PasswordInput } from './PasswordInput';  // Make sure to update with correct import path
import '@testing-library/jest-dom';

describe('PasswordInput', () => {
    const mockOnChange = jest.fn(event => {
        return {
            name: event.target.name,
            value: event.target.value
        }
    });

    it('renders a password input field with a label', () => {
        render(<PasswordInput value="" onChange={mockOnChange} />);
        const input = screen.getByTestId('password-input');
        expect(input).toBeInTheDocument();
    });

    it('calls the onChange callback with the input value', () => {
        render(<PasswordInput value="" onChange={mockOnChange} />);
        const passwordInput = screen.getByTestId('password-input') as HTMLInputElement;
        fireEvent.change(passwordInput, { target: { value: 'testpassword' } });
        expect(mockOnChange).toHaveBeenCalled();
        expect(mockOnChange).toReturnWith({name: 'password', value: 'testpassword'});
    });

    it('renders a password input field with the correct name attribute', () => {
        render(<PasswordInput value="" onChange={mockOnChange} />);
        const passwordInput = screen.getByTestId('password-input') as HTMLInputElement;
        expect(passwordInput.name).toBe('password');
    });

    it('toggles password visibility', () => {
        render(<PasswordInput value="" onChange={mockOnChange} />);
        const passwordInput = screen.getByTestId('password-input') as HTMLInputElement;
        const toggleButton = screen.getByRole('button');

        // The password field should initially be hidden
        expect(passwordInput.type).toBe('password');

        // Clicking the button should reveal the password
        fireEvent.click(toggleButton);
        expect(passwordInput.type).toBe('text');

        // Clicking the button again should hide the password
        fireEvent.click(toggleButton);
        expect(passwordInput.type).toBe('password');
    });
});
