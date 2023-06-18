import {
    Button,
    Stack,
} from '@chakra-ui/react';
import { useRouter } from 'next/router';

const FormButtons: React.FC<{handleSubmit: () => void}> = ({handleSubmit}) => {
    const router = useRouter();
    const handleSkip = () => {
        router.push(`/profile/me`)
    }

    return (
        <Stack isInline justifyContent="flex-end">
            <Button variant="outline" mr={2} onClick={handleSkip}>
                Skip
            </Button>
            <Button colorScheme="teal" onClick={handleSubmit}>
                Save
            </Button>
        </Stack>
    );
}

export default FormButtons;