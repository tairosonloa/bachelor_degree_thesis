import React from "react"
import styles from "./cards.module.css"


class Cards extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      reservations: []
    }
    this.currentHour = 0
    this.currentMinutes = 0
    try { // Load config
      this.config = require("/etc/rpi3_conf.json")
    } catch {
      this.config = require("../../../../../config.json")
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
   * Updates the current time, used to display only future reservations
   */
  updateCurrentTime = () => {
    let date = new Date();
    this.currentHour = date.getHours();
    this.currentMinutes = date.getMinutes();
  }

  /**
   * Returns the html <div> object of a reservation
   * @param {reservation dict}  r dictionary with all info of a reservation
   * @param {int}               i index to set as react key of the div
   */
  getCard = (r, i) => {
    if (this.currentHour < r["EndHour"] || (this.currentHour == r["EndHour"] && this.currentMinutes < r["EndMinute"])) {
      return <div key={i} className={styles.card}>
        <div className={styles.subject}>{r["Subject"]}</div>
        <div className={styles.study}>{r["Study"]}</div>
        <div className={styles.classroom}>{r["Classroom"]} de {r["StartHour"] + ":" + 
          (r["StartMinute"] === 0 ? "00" : r["StartMinute"])} a {r["EndHour"] + ":" + 
          (r["EndMinute"] === 0 ? "00" : r["EndMinute"])}</div>
        <div className={styles.professor}>{r["Professor"]}</div>
      </div>;
    }
    return null
  }

  /**
   * Returns an array of html <div>, where every <div> is a reservation card
   */
  createCards = () => {
    this.updateCurrentTime()
    let cards = []
    for (const [i, r] of this.state.reservations.entries()) {
      let card = this.getCard(r, i)
      if (card != null) {
        cards.push(card)
      }
      if (cards.length === 4) {
        break;
      }
    }
    if (cards.length !== 0) {
      return cards
    }
    return <div className={styles.endCard}>No hay reservas para el día de hoy o ya han finalizado todas las reservas</div>
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
    }, 60000);
  }
  componentWillUnmount() {
    clearInterval(this.timer);
  }
}

export default Cards;