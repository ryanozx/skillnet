import React from 'react';
import { screen, render, fireEvent, getByTestId } from '@testing-library/react';
import { EmailInput } from './EmailInput';
import '@testing-library/jest-dom';

describe('EmailInput', () => {
    const mockOnChange = jest.fn(event => {
        return {
            name: event.target.name,
            value: event.target.value
        }
    });

    it('renders an email input field with a label', () => {
        render(<EmailInput value="" onChange={mockOnChange} />);
        const input = screen.getByTestId('email-input');
        expect(input).toBeInTheDocument();
    });

    it('renders an email input field with the correct name attribute', () => {
        render(<EmailInput value="" onChange={mockOnChange} />);
        const emailInput = screen.getByTestId('email-input') as HTMLInputElement;
        expect(emailInput.name).toBe('email');
    });

    it('calls the onChange callback with the input value', () => {
        render(<EmailInput value="" onChange={mockOnChange} />);
        const emailInput = screen.getByTestId('email-input') as HTMLInputElement;
        fireEvent.change(emailInput, { target: { value: 'test@email.com' } });
        expect(mockOnChange).toHaveBeenCalled();
        expect(mockOnChange).toReturnWith({name: 'email', value: 'test@email.com'});
    });    
});



