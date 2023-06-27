// import React from 'react';
// import { render, fireEvent, waitFor } from '@testing-library/react';
// import axios from 'axios';
// import { useRouter } from 'next/router';
// import { useToast } from '@chakra-ui/react';
// import SignUpForm from './SignUpForm';

// jest.mock('axios');
// jest.mock('next/router', () => ({
//     useRouter: jest.fn(),
// }));

// jest.mock('@chakra-ui/react', () => ({
//     useToast: jest.fn(),
// }));

// describe('SignUpForm', () => {
//     it('calls the signup API and redirects on successful submission', async () => {
//         const router = {
//             push: jest.fn(),
//         };
//         const toast = jest.fn();
//         (useRouter as jest.Mock).mockReturnValue(router);
//         (useToast as jest.Mock).mockReturnValue(toast);
//         (axios.post as jest.Mock).mockResolvedValue({});

//         const { getByTestId } = render(<SignUpForm />);
//         fireEvent.change(getByTestId('username-input'), { target: { value: 'test' } });
//         fireEvent.change(getByTestId('email-input'), { target: { value: 'test@example.com' } });
//         fireEvent.change(getByTestId('password-input'), { target: { value: 'password' } });
//         fireEvent.click(getByTestId('submit-button'));

//         await waitFor(() => expect(axios.post).toHaveBeenCalled());
//         expect(router.push).toHaveBeenCalledWith('/create-profile');
//         expect(toast).toHaveBeenCalledWith({
//             title: "Form submission successful.",
//             description: "Account successfully created.",
//             status: "success",
//             duration: 5000,
//             isClosable: true,
//         });
//     });
// });
