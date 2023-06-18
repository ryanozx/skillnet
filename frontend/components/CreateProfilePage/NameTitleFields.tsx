import React from 'react';
import {
    FormControl,
    FormLabel,
    Input,
    Stack,
} from '@chakra-ui/react';

const NameTitleFields: React.FC<{form: any, handleChange: (e: React.ChangeEvent<HTMLInputElement>) => void}> = ({form, handleChange}) => {
    return (
        <>
            <Stack>
                <FormControl>
                    <FormLabel>Name</FormLabel>
                    <Input
                        type="text"
                        name="name"
                        value={form.name}
                        onChange={handleChange}
                    />
                </FormControl>
                <FormControl>
                    <FormLabel>Title</FormLabel>
                    <Input
                        type="text"
                        name="title"
                        value={form.title}
                        onChange={handleChange}
                    />
                </FormControl>
            </Stack>
        </>
    );
}

export default NameTitleFields;
