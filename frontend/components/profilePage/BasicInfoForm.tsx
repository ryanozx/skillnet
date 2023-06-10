import React from 'react';
import {
    ModalHeader,
    ModalBody,
    ModalFooter,
    FormControl,
    FormLabel,
    Input,
    Textarea,
    IconButton,
    Button,
    Switch
} from '@chakra-ui/react';


export default function BasicInfoForm(props: any) {
    const { form, handleInputChange } = props;
    return (
        <>
            <ModalHeader>Edit your profile</ModalHeader>
            <ModalBody>
                <FormControl id="name">
                    <FormLabel>Name</FormLabel>
                    <Input name="name" value={form.name} onChange={handleInputChange} placeholder="Your name" />
                </FormControl>
                <FormControl id="title" mt={4}>
                    <FormLabel>Title</FormLabel>
                    <Input name="title" value={form.title} onChange={handleInputChange} placeholder="Your title" />
                </FormControl>
                <FormControl id="about" mt={4}>
                    <FormLabel>About me</FormLabel>
                    <Textarea name="about" value={form.about} onChange={handleInputChange} placeholder="About you" />
                </FormControl>
            </ModalBody>
        </>
    );
}
