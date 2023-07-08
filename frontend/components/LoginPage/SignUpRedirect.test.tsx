import React from 'react';
import { render } from '@testing-library/react';
import SignUpRedirect from './SignUpRedirect';

describe('SignUpRedirect', () => {
    it('renders the redirect text', () => {
        const { getByTestId } = render(<SignUpRedirect />);
        const redirectText = getByTestId('redirect-text');
        
        expect(redirectText).toBeInTheDocument();
        expect(redirectText.textContent).toContain("Don't have an account? Sign up");
    });

    it('renders the login link', () => {
        const { getByTestId } = render(<SignUpRedirect />);
        const loginLink = getByTestId('signup-link');
        
        expect(loginLink).toBeInTheDocument();
        expect(loginLink.getAttribute('href')).toBe('/signup');
    });
});
