import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import { SubmitButton } from './SubmitButton'; // adjust the path as necessary

describe('SubmitButton', () => {
    it('renders a submit button', () => {
        const mockOnClick = jest.fn();
        const { getByText } = render(<SubmitButton onClick={mockOnClick} />);

        const button = getByText('Sign up');
        expect(button).toBeInTheDocument();
    });

    it('calls the onClick handler when clicked', () => {
        const mockOnClick = jest.fn();
        const { getByText } = render(<SubmitButton onClick={mockOnClick} />);

        fireEvent.click(getByText('Sign up'));

        expect(mockOnClick).toHaveBeenCalled();
    });
});
