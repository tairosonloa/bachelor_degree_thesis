import React from "react"

class Reservations extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      reservations: []
    }
    try { // Load config
      this.config = require("/etc/rpi3_conf.json")
    } catch {
      this.config = require("../../../../config.json")
    }
  }

  /**
   * Makes HTTP GET request to rpi3 API to get JSON including today reservations
   */
  getReservations = () => {
    fetch("http://" + this.config.Rpi3APIAddress + ":" + this.config.Rpi3APIPort + "/reservations")
      .then(response => response.json())
      .then(json => this.setState({reservations: json}))
      .catch(error => console.log('Request failed', error))
  }

  /**
   * Returns the html <div> object of a reservation
   * @param {reservation dict}  r dictionary with all info of a reservation
   * @param {int}               i index to set as react key of the div
   */
  getCard = (r, i) => {
    return <div key={i}><div>Asignatura: {r["Subject"]}</div><div>{(r["Study"] === "" ? "" : r["Study"])}</div><div>Aula: {r["Classroom"]} de {r["StartTime"]} a {r["EndTime"]}</div><div>{(r["Professor"] === "" ? "" : r["Professor"])}</div><br/></div>
  }

  /**
   * Returns an array of html <div>, where every <div> is a reservation card
   */
  createCards = () => {
    let cards = []
    let i = 0
    for (let r of this.state.reservations) {
      cards.push(this.getCard(r, i))
      i++
    }
    return cards
  }

  render() {
    return ( // TODO: maybe <tag>{function()}</tag>
      this.createCards()
    );
  }

  componentDidMount() {
    this.getReservations()
    this.timer = setInterval(() => {
      this.getReservations()
    }, 1800000);
  }
  componentWillUnmount() {
    clearInterval(this.timer);
  }
}

export default Reservations;