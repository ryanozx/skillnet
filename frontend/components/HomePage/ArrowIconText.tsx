import React from 'react';
import { createIcon, Box, Icon, Text, useColorModeValue } from '@chakra-ui/react';
import { ArrowIcon } from './ArrowIcon';
export const ArrowIconText = () => (
    <>
        <Box>
            <Icon
                as={ArrowIcon}
                color={useColorModeValue('gray.800', 'gray.300')}
                w={71}
                position={'absolute'}
                right={-71}
                top={'10px'}
            />
            <Text
                fontSize={'lg'}
                fontFamily={'Caveat'}
                position={'absolute'}
                right={'-125px'}
                top={'-15px'}
                transform={'rotate(10deg)'}>
                Join us now!
            </Text>
        </Box>
    </>
);

