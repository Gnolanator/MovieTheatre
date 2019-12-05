import React, { Component, useEffect, useState} from 'react';
import logo from './logo.svg';
import './App.css';
import MovieCarousel from "./components/Movies";
import UploadButton from "./components/UploadButton";


class App extends Component {
  render() {
    return (
      <div className="App">
        <UploadButton />
      </div>
    );
  }
}

export default App;
