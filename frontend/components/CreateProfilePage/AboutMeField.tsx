import React from 'react';
import {
    FormControl,
    FormLabel,
    Textarea,
} from '@chakra-ui/react';

const AboutMeField: React.FC<{form: any, handleChange: (e: React.ChangeEvent<HTMLTextAreaElement>) => void}> = ({form, handleChange}) => {
    return (
        <FormControl>
            <FormLabel>About Me</FormLabel>
            <Textarea
                name="aboutMe"
                value={form.aboutMe}
                onChange={handleChange}
            />
        </FormControl>
    );
}

export default AboutMeField;