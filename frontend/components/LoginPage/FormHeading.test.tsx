import React from 'react';
import { render } from '@testing-library/react';
import FormHeading from './FormHeading';

describe('FormHeading', () => {
    it('renders the main heading', () => {
        const { getByTestId } = render(<FormHeading />);
        const heading = getByTestId('form-heading');
        
        expect(heading).toBeInTheDocument();
        expect(heading.textContent).toBe('Sign in to your account');
    });

    it('renders the subheading', () => {
        const { getByTestId } = render(<FormHeading />);
        const subheading = getByTestId('form-subheading');
        
        expect(subheading).toBeInTheDocument();
        expect(subheading.textContent).toBe('to enjoy all of our cool features ✌️');
    });
});
