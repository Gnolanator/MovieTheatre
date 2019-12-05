import React, { useState, useEffect } from 'react';

const getMovies = () => {
  const [movies, setMovies] = useState(0);

  useEffect(() => {
    fetch('/movies').then(response =>
      response.json().then(data => {
        console.log(data);
        setMovies(data);
      })
    );
  }

