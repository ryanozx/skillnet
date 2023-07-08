import React from 'react';
import { render } from '@testing-library/react';
import HomePageHeader from './HomePageHeader';

describe('HomePageHeader', () => {
    it('renders the correct header', () => {
        const { getByTestId } = render(<HomePageHeader />);
        const header = getByTestId('home-page-header');
        expect(header).toBeInTheDocument();
    });

    it('renders the correct subheader', () => {
        const { getByTestId } = render(<HomePageHeader />);
        const subheader = getByTestId('home-page-subheader');
        expect(subheader).toBeInTheDocument();
    });
});
