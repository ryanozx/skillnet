import { Box, Stack, FormControl, FormLabel, Input, Button, useColorModeValue, useToast } from '@chakra-ui/react';
import ActionSection from './ActionSection';
import React, { useState, MouseEventHandler } from "react";
import axios from "axios";
import { useRouter } from 'next/router'; 

type User = {
    username: string;
    password: string;
};

export default function LoginForm() {
    const [form, setForm] = useState<User>({ username: "", password: "" });
    const toast = useToast();
    const router = useRouter();

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setForm((prevForm) => ({ ...prevForm, [name]: value }));
    };
    const onSubmit: MouseEventHandler = () => {
        const {username, password} = form
        var form_data = new FormData();
        form_data.append('username', username);
        form_data.append('password', password);
        axios.post('http://localhost:8080/login', form_data, {withCredentials: true})
            .then((res) => {
                console.log(res);
                toast({
                    title: "Form submission successful.",
                    description: "We've received your form data.",
                    status: "success",
                    duration: 5000,
                    isClosable: true,
                });
                router.push("/profile/me"); // Navigate to "/"
            })
            .catch((error) => {
                console.log(error);
                toast({
                    title: "An error occurred.",
                    description: error.response.data.error,
                    status: "error",
                    duration: 5000,
                    isClosable: true,
                });
            });
    };

    return (
        <Box
            rounded={'lg'}
            bg={useColorModeValue('white', 'gray.700')}
            boxShadow={'lg'}
            p={8}
            w={{ base: '90vw', md: '60vw', lg: '30vw' }}
        >
            <Stack spacing={4}>
                <FormControl id="username">
                    <FormLabel>username</FormLabel>
                    <Input
                        type="username"
                        name="username"
                        value={form.username}
                        onChange={handleInputChange}
                    />
                </FormControl>
                <FormControl id="password">
                    <FormLabel>Password</FormLabel>
                    <Input
                        type="password"
                        name="password"
                        value={form.password}
                        onChange={handleInputChange}
                    />
                </FormControl>
                <ActionSection onSubmit={onSubmit}/>
            </Stack>
        </Box>
    );
}

// import React, { MouseEventHandler, useState } from 'react';
// import axios from 'axios';
// import { useToast } from "@chakra-ui/react";
// import { Box, Stack, FormControl, FormLabel, Input } from "@chakra-ui/react";

// // Custom hook for form management
// const useForm = (initialState: any) => {
//   const [form, setForm] = useState(initialState);

//   const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
//     const { name, value } = e.target;
//     setForm((prevForm) => ({ ...prevForm, [name]: value }));
//   };

//   return { form, handleInputChange };
// };

// // Custom hook for login
// const useLogin = (form: any) => {
//   const toast = useToast();

//   const onSubmit: MouseEventHandler = () => {
//     const { username, password } = form;
//     var form_data = new FormData();
//     form_data.append('username', username);
//     form_data.append('password', password);
//     axios
//       .post('http://localhost:8080/login', form_data, { withCredentials: true })
//       .then((res) => {
//         console.log(res);
//         toast({
//           title: "Form submission successful.",
//           description: "We've received your form data.",
//           status: "success",
//           duration: 5000,
//           isClosable: true,
//         });
//       })
//       .catch((error) => {
//         console.log(error);
//         toast({
//           title: "An error occurred.",
//           description: error.response.data.error,
//           status: "error",
//           duration: 5000,
//           isClosable: true,
//         });
//       });
//   };

//   return { onSubmit };
// };

// interface User {
//   username: string;
//   password: string;
// }

// export default function LoginForm() {
//   const { form, handleInputChange } = useForm<User>({ username: "", password: "" });
//   const { onSubmit } = useLogin(form);

//   return (
//     <Box
//       rounded={'lg'}
//       bg={useColorModeValue('white', 'gray.700')}
//       boxShadow={'lg'}
//       p={8}
//       w={{ base: '90vw', md: '60vw', lg: '30vw' }}
//     >
//       <Stack spacing={4}>
//         <FormControl id="username">
//           <FormLabel>username</FormLabel>
//           <Input
//             type="username"
//             name="username"
//             value={form.username}
//             onChange={handleInputChange}
//           />
//         </FormControl>
//         <FormControl id="password">
//           <FormLabel>Password</FormLabel>
//           <Input
//             type="password"
//             name="password"
//             value={form.password}
//             onChange={handleInputChange}
//           />
//         </FormControl>
//         <ActionSection onSubmit={onSubmit} />
//       </Stack>
//     </Box>
//   );
// }
