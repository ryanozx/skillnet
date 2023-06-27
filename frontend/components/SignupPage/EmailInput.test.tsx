import React from 'react';
import { render, fireEvent, getByTestId } from '@testing-library/react';
import { EmailInput } from './EmailInput';
import '@testing-library/jest-dom';


describe('EmailInput', () => {
    it('renders an email input field with a label', () => {
        const mockOnChange = jest.fn();
        const { getByTestId } = render(<EmailInput value="" onChange={mockOnChange} />);
        const input = getByTestId('email-input');
        expect(input).toBeInTheDocument();
    });

    it('calls the onChange handler when typing', () => {
        const mockOnChange = jest.fn();
        const { getByTestId } = render(<EmailInput value="" onChange={mockOnChange} />);
        const input = getByTestId('email-input');
        fireEvent.change(input, { target: { value: 'test@test.com' } });

        expect(mockOnChange).toHaveBeenCalledWith('test@test.com');
    });
});
