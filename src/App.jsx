import logo from './logo.svg';
import './App.css';
import axios from 'axios'
import MainPage from './MainPage'
import Signup from './Authentication/Signup'
import HomePage from './Authentication/HomePage'
import React, {useState} from "react";
import {
  BrowserRouter as Router,
  Route,
  Routes,
  Link,
} from "react-router-dom"; 

const App =()=>{
 
 return(
  <>
  <Router>
  <Routes>
  <Route path="/account" element={<HomePage/>}></Route>
  <Route path="/" element={<Signup/>}>
  </Route>
  </Routes>
  </Router>
  </>
    );
}

export default App;