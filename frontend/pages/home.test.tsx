import React from 'react';
import { render } from '@testing-library/react';
import Home from './index';

jest.mock('next/router', () => require('next-router-mock'));
jest.mock('../withAuthRedirect', () => ({
    preventAuthAccess: jest.fn((Component) => Component),
    requireAuth: jest.fn((Component) => Component),
}));

it('Home page renders correctly', () => {
    const { asFragment } = render(<Home />);
    expect(asFragment()).toMatchSnapshot();
});
