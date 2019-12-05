import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import "react-responsive-carousel/lib/styles/carousel.min.css";
import { Carousel } from 'react-responsive-carousel';

export default class MovieCarousel extends Component {

  state = {
    file: null,
    movies: []
  }

  componentDidMount() {
    fetch('/movies').then(response =>
      response.json().then(data => {
        console.log(data);
        this.setState({ movies: data });
      })
    );
  }

  handleFile(e) {
    let file = e.target.files[0]
    this.setState({file: file})
  }

  handleUpload(e) {
    console.log(this.state.file, "THE STATE");
  }

  render () {
    console.log(this.state.movies);
    return (
      <div className='Button'>
        <h1>Current Catalog:</h1>
        <Carousel>
          {this.state.movies.map(movie =>
            <div>
              <img src={movie[2]} />
              <p className="legend">{movie[1]}</p>
            </div>
          )}
        </Carousel>
        <h1>UPLOAD</h1>
        <form>
          <div className='select'>
            <label>Select File: </label>
            <input type='file' name='file' onChange={(e) => this.handleFile(e)} />
          </div>

          <button type='button' onClick={(e) => this.handleUpload(e)}>Upload</button>
        </form>
      </div>
    );
  }
}
