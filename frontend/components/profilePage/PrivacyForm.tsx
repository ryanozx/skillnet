import {
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalCloseButton,
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

export default function PrivacyForm(props: any) {
    const { handleSwitchChange, form } = props;
    return (
        <>
            <ModalHeader>Privacy Settings</ModalHeader>
            <ModalBody>
                <FormControl display="flex" alignItems="center">
                    <FormLabel mb="0">Tagline</FormLabel>
                    <Switch size="lg" ml="auto" isChecked={form.privacySettings.tagline} onChange={() => handleSwitchChange("tagline")} />
                </FormControl>
                <FormControl display="flex" alignItems="center" mt={4}>
                    <FormLabel mb="0">About me</FormLabel>
                    <Switch size="lg" ml="auto" isChecked={form.privacySettings.about} onChange={() => handleSwitchChange("about")}/>
                </FormControl>
                <FormControl display="flex" alignItems="center" mt={4}>
                    <FormLabel mb="0">Projects</FormLabel>
                    <Switch size="lg" ml="auto" isChecked={form.privacySettings.projects} onChange={() => handleSwitchChange("projects")}/>
                </FormControl>
                <FormControl display="flex" alignItems="center" mt={4}>
                    <FormLabel mb="0">Activity</FormLabel>
                    <Switch size="lg" ml="auto" isChecked={form.privacySettings.activity} onChange={() => handleSwitchChange("activity")}/>
                </FormControl>
            </ModalBody>
        </>

    );

}