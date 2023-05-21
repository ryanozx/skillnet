import React, { useRef } from 'react';
import { Avatar, IconButton, Input } from '@chakra-ui/react';

export default function EditPicBtn({ currentProfilePic }) {
  const fileInputRef = useRef();

  const handleFileInputClick = () => {
    fileInputRef.current.click();
  };

  const handleFileChange = (event) => {
    const file = event.target.files[0];
    // handle the selected file, e.g. upload it
  };

  return (
    <>
        <IconButton
            aria-label="Change profile picture"
            icon={<Avatar size="2xl" src={currentProfilePic} />}
            onClick={handleFileInputClick}
            isRound="true"
        />

        <Input
            type="file"
            accept="image/*"
            ref={fileInputRef}
            onChange={handleFileChange}
            style={{ display: 'none' }}
        />
    </>
  );
};

