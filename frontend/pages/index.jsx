// import Head from 'next/head';
// import styles from '../styles/Home.module.css';
import {
  Accordion,
  AccordionItem,
  AccordionButton,
  AccordionPanel,
  AccordionIcon,
  Box
} from '@chakra-ui/react'
import React from 'react';
import NavBar from '../components/base/NavBar';
import Searchbar from '../components/base/Searchbar';
import SideBar from '../components/base/SideBar';

export default function Home() {
  return (
    <>
      <NavBar></NavBar>
      <SideBar></SideBar>      
    </>
  )
}
