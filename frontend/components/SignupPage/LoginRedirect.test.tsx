import React from 'react';
import { render } from '@testing-library/react';
import LoginRedirect from './LoginRedirect';

describe('LoginRedirect', () => {
    it('renders the redirect text', () => {
        const { getByTestId } = render(<LoginRedirect />);
        const redirectText = getByTestId('redirect-text');
        
        expect(redirectText).toBeInTheDocument();
        expect(redirectText.textContent).toContain('Already a user?');
    });

    it('renders the login link', () => {
        const { getByTestId } = render(<LoginRedirect />);
        const loginLink = getByTestId('login-link');
        
        expect(loginLink).toBeInTheDocument();
        expect(loginLink.getAttribute('href')).toBe('/login');
    });
});
