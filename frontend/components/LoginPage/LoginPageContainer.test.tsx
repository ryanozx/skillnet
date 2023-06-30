import React from 'react';
import { render } from '@testing-library/react';
import LoginPageContainer from './LoginPageContainer';
import { preventAuthAccess } from '../../withAuthRedirect';

// Mock the components to isolate the container component testing
jest.mock('./FormHeading', () => () => <div data-testid="form-heading" />);
jest.mock('./LoginForm', () => () => <div data-testid="login-form" />);
jest.mock('../../withAuthRedirect', () => ({
    preventAuthAccess: jest.fn((Component) => Component),
}));

describe('LoginPageContainer', () => {
    it('renders the FormHeading and LoginForm components', () => {
        const { getByTestId } = render(<LoginPageContainer />);
        
        expect(getByTestId('form-heading')).toBeInTheDocument();
        expect(getByTestId('login-form')).toBeInTheDocument();
    });

    it('calls preventAuthAccess to prevent authenticated user from visiting this page', () => {
        expect(preventAuthAccess).toHaveBeenCalled();
    });

});

