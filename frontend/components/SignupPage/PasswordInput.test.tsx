import React from 'react';
import { render, fireEvent, getByTestId } from '@testing-library/react';
import { PasswordInput } from './PasswordInput';  // Make sure to update with correct import path
import '@testing-library/jest-dom';

describe("PasswordInput", () => {
  it('renders a password input field initially', () => {
    const { getByLabelText } = render(<PasswordInput value="" onChange={() => {}} />);
    const passwordInput = getByLabelText("Password") as HTMLInputElement;
    expect(passwordInput.type).toBe('password');
  });

  it('toggles input type between password and text when show/hide button is clicked', () => {
    const { getByLabelText, getByRole } = render(<PasswordInput value="" onChange={() => {}} />);
    const passwordInput = getByLabelText("Password") as HTMLInputElement;
    const toggleButton = getByRole('button');

    fireEvent.click(toggleButton);
    expect(passwordInput.type).toBe('text');

    fireEvent.click(toggleButton);
    expect(passwordInput.type).toBe('password');
  });

  it('calls the onChange function with input value when it changes', () => {
    const handleChange = jest.fn();
    const { getByLabelText } = render(<PasswordInput value="" onChange={handleChange} />);
    const passwordInput = getByLabelText("Password");

    fireEvent.change(passwordInput, { target: { value: 'New Password' } });
    expect(handleChange).toHaveBeenCalledWith('New Password');
  });
});
