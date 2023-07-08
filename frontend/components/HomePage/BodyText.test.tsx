import React from 'react';
import { render } from '@testing-library/react';
import { BodyText } from './BodyText';

describe('BodyText', () => {

    it('renders the correct text', () => {
        const { getByTestId } = render(<BodyText />);
        const bodyText = getByTestId('body-text-content');
        expect(bodyText).toBeInTheDocument();
    });
});
