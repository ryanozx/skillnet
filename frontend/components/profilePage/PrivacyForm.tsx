import {
    ModalHeader,
    ModalBody,
    FormControl,
    FormLabel,
    Switch
} from '@chakra-ui/react';
import { FormType } from './EditInfoModal';

interface PrivacyFormProps {
    handleSwitchChange: (name : string) => void
    form: FormType
}

export default function PrivacyForm(props: PrivacyFormProps) {
    const { handleSwitchChange, form } = props;
    return (
        <>
            <ModalHeader>Privacy Settings</ModalHeader>
            <ModalBody>
                <FormControl display="flex" alignItems="center">
                    <FormLabel mb="0">Title</FormLabel>
                    <Switch size="lg" ml="auto" isChecked={form.privacySettings.title} onChange={() => handleSwitchChange("title")} />
                </FormControl>
                <FormControl display="flex" alignItems="center" mt={4}>
                    <FormLabel mb="0">About me</FormLabel>
                    <Switch size="lg" ml="auto" isChecked={form.privacySettings.about} onChange={() => handleSwitchChange("about")}/>
                </FormControl>
                {/* <FormControl display="flex" alignItems="center" mt={4}>
                    <FormLabel mb="0">Projects</FormLabel>
                    <Switch size="lg" ml="auto" isChecked={form.privacySettings.projects} onChange={() => handleSwitchChange("projects")}/>
                </FormControl>
                <FormControl display="flex" alignItems="center" mt={4}>
                    <FormLabel mb="0">Activity</FormLabel>
                    <Switch size="lg" ml="auto" isChecked={form.privacySettings.activity} onChange={() => handleSwitchChange("activity")}/>
                </FormControl> */}
            </ModalBody>
        </>

    );

}