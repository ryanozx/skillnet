import React from 'react';
import { render } from '@testing-library/react';
import { preventAuthAccess } from '../../withAuthRedirect';
import HomePageContainer from './HomePageContainer';


jest.mock('./HomePageHeader', () => () => <div data-testid="home-page-header" />);
jest.mock('./CallToActionButtons', () => () => <div data-testid="call-to-action-buttons" />);
jest.mock('../../withAuthRedirect', () => ({
    preventAuthAccess: jest.fn((Component) => Component),
    requireAuth: jest.fn((Component) => Component),
}));

describe('HomePageContainer', () => {
    it('renders the HomePageHeader and CallToActionButtons components', () => {
        const { getByTestId } = render(<HomePageContainer />);
        
        expect(getByTestId('home-page-header')).toBeInTheDocument();
        expect(getByTestId('call-to-action-buttons')).toBeInTheDocument();
    });

    it('calls preventAuthAccess to prevent authenticated user from visiting this page', () => {
        expect(preventAuthAccess).toHaveBeenCalled();
    });
});
