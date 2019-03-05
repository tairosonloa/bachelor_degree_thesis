import React from "react"
import styles from "./classrooms.module.css"


class Classrooms extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classrooms: []
    }
    try { // Load config
      this.config = require("/etc/rpi3_conf.json")
    } catch {
      this.config = require("../../../../../config.json")
    }
  }

  /**
   * Makes HTTP GET request to rpi3 API to get JSON classrooms status
   */
  getClassrooms = () => {
    fetch("http://" + this.config.Rpi3APIAddress + ":" + this.config.Rpi3APIPort + "/classrooms")
      .then(response => response.json())
      .then(json => this.setState({classrooms: json}))
      .catch(error => console.log('Request failed', error))
  }


  /**
   * Returns an array of html <div>, where every <div> is a classroom
   */
  createClassrooms = () => {
    let divs = []
    let key = 0
    for (let c in this.state.classrooms) {
      if (this.state.classrooms[c] === 0 ) {
        divs.push(<div key={key} className={styles.free}>{c}</div>)
      } else if (this.state.classrooms[c] === 1 ) {
        divs.push(<div key={key} className={styles.occupied}>{c}</div>)
      } else if (this.state.classrooms[c] === 2 ) {
        divs.push(<div key={key} className={styles.reserved}>{c}</div>)
      } else {
        divs.push(<div key={key} className={styles.futureOccupied}>{c}</div>)
      }
      key++
    }
    return divs
  }

  render() {
    return ( // TODO: maybe <tag>{function()}</tag>
      this.createClassrooms()
    );
  }

  componentDidMount() {
    this.getClassrooms()
    this.timer = setInterval(() => {
      this.getClassrooms()
    }, 60000);
  }
  componentWillUnmount() {
    clearInterval(this.timer);
  }
}

export default Classrooms;