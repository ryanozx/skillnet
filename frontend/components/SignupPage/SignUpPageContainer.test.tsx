import React from 'react';
import { render } from '@testing-library/react';
import SignUpPageContainer from './SignUpPageContainer';
import { preventAuthAccess } from '../../withAuthRedirect';

jest.mock('./FormHeading', () => () => <div data-testid="form-heading" />);
jest.mock('./SignUpForm', () => () => <div data-testid="sign-up-form" />);
jest.mock('../../withAuthRedirect', () => ({
    preventAuthAccess: jest.fn((Component) => Component),
    requireAuth: jest.fn((Component) => Component),
}));
  
describe('SignUpPageContainer', () => {
    it('renders the FormHeading and SignUpForm components', () => {
        const { getByTestId } = render(<SignUpPageContainer />);
        
        expect(getByTestId('form-heading')).toBeInTheDocument();
        expect(getByTestId('sign-up-form')).toBeInTheDocument();
    });

    it('calls preventAuthAccess to prevent authenticated user from visiting this page', () => {
        expect(preventAuthAccess).toHaveBeenCalled();
    });
});
