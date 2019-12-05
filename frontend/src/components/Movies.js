import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import "react-responsive-carousel/lib/styles/carousel.min.css";
import { Carousel } from 'react-responsive-carousel';

class MovieCarousel extends Component {
  state = {}

  render () {
    console.log(this.props.movies);
    return (
      <Carousel>
      </Carousel>
    );
  }
};

export default MovieCarousel;
