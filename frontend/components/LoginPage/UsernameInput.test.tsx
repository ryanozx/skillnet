import React from 'react';
import { screen, render, cleanup, fireEvent } from '@testing-library/react';
import { UsernameInput } from './UsernameInput';

describe('UsernameInput', () => {
    const mockOnChange = jest.fn(event => {
        return {
            name: event.target.name,
            value: event.target.value
        }
    });

    it('renders a username input field with a label', () => {
        render(<UsernameInput value="" onChange={mockOnChange} />);
        const input = screen.getByTestId('username-input');
        expect(input).toBeInTheDocument();
    });

    it('calls the onChange callback with the input value', () => {
        render(<UsernameInput value="" onChange={mockOnChange} />);
        const usernameInput = screen.getByTestId('username-input') as HTMLInputElement;
        fireEvent.change(usernameInput, { target: { value: 'testusername' } });
        expect(mockOnChange).toHaveBeenCalled();
        expect(mockOnChange).toReturnWith({name: 'username', value: 'testusername'});
    });

    it('renders a username input field with the correct name attribute', () => {
        render(<UsernameInput value="" onChange={mockOnChange} />);
        const usernameInput = screen.getByTestId('username-input') as HTMLInputElement;
        expect(usernameInput.name).toBe('username');
    });
});
