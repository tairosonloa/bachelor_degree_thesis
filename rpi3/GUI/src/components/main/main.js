import React from "react"
import styles from "./main.module.css"


class Main extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      reservations: [],
      occupation : []
    }
    this.rotated = false
    this.reservationsNum = 0
    this.currentHour = 0
    this.currentMinutes = 0
    this.globalState = 0
    this.classroomToShow = 0
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
    if (this.rotated && i < 4 && this.reservationsNum > 4) {
      // if rotating, show next reservations (i >= 4)
      return null
    }
    if (this.currentHour < r["EndHour"] || (this.currentHour === r["EndHour"] && this.currentMinutes < r["EndMinute"])) {
      return <div key={i} className={styles.card}>
        <div className={styles.subject}>{r["Subject"]}</div>
        <div className={styles.study}>{r["Study"]}</div>
        <div className={styles.classroom}>{r["Classroom"]} de {r["StartHour"] + ":" + 
          (r["StartMinute"] === 0 ? "00" : r["StartMinute"])} a {r["EndHour"] + ":" + 
          (r["EndMinute"] === 0 ? "00" : r["EndMinute"])}</div>
      </div>;
    }
    return null
  }

  /**
   * Returns an array of html <div>, where every <div> is a reservation card
   */
  createCards = () => {
    this.updateCurrentTime()
    let cards = [];
    // Get reservations to show
    for (const [i, r] of this.state.reservations.entries()) {
      let card = this.getCard(r, i)
      if (card != null) {
        cards.push(card)
      }
      if (cards.length === 4) {
        break;
      }
    }
    // Change between rotation states and update global state
    (this.rotated) ? this.rotated = false : this.rotated = true
    this.globalState = (this.globalState + 1) % 6
    // Return cards if any
    if (cards.length !== 0) {
      this.reservationsNum = this.state.reservations.length
      return cards
    }
    return <div className={styles.endCard}>No hay reservas para el día de hoy o ya han finalizado todas las reservas</div>
  }

  getOccupation = () => {
    fetch("http://" + this.config.Rpi3APIAddress + ":" + this.config.Rpi3APIPort + "/occupation")
      .then(response => response.json())
      .then(json => {
        var classrooms = ["caca"]
        console.log(classrooms)
        classrooms[0] = json["F16"].Computers
        classrooms[1] = json["F18"].Computers
        classrooms[2] = json["C05"].Computers
        classrooms[3] = json["C06"].Computers
        console.log(classrooms)
        this.setState({occupation: classrooms})
      })
      .catch(error => console.log('Request failed', error))
  }

  getComputer = (r, i, name) => {
    let ip = i + 1
    if (name.includes("C")) {
      ip += 50
    }
    switch (r) {
      case 0:
        return <div key={i} className={styles.shutdown}>{(ip < 10)? name + "0"+ ip.toString() : name + ip.toString()}</div>
      case 1:
        return <div key={i} className={styles.linux}>{(ip < 10)? name + "0"+ ip.toString() : name + ip.toString()}</div>
      case 2:
        return <div key={i} className={styles.windows}>{(ip < 10)? name + "0"+ ip.toString() : name + ip.toString()}</div>
      case 3:
        return <div key={i} className={styles.linuxUser}>{(ip < 10)? name + "0"+ ip.toString() : name + ip.toString()}</div>
      case 4:
        return <div key={i} className={styles.windowsUser}>{(ip < 10)? name + "0"+ ip.toString() : name + ip.toString()}</div>
      case 5:
        return <div key={i} className={styles.timeout}>{(ip < 10)? name + "0"+ ip.toString() : name + ip.toString()}</div>
      default:
        return <div key={i}></div>
    }
  }

  printClassrooms = () => {
    let computers = ["F1", "F2", "C", "C1"]
    let classroom = ["4.0.F16", "4.0.F18", "2.2.C05", "2.2.C06"]
    let classroomMap = [<h2 className={styles.title}>Aula {classroom[this.classroomToShow]}</h2>]
    // Get computer status of the classroom
    console.log(this.classroomToShow)
    console.log(this.state.occupation)
    console.log(this.state.occupation[this.classroomToShow])
    for (const [i, r] of this.state.occupation[this.classroomToShow].entries()) {
      classroomMap.push(this.getComputer(r, i, computers[this.classroomToShow]))
    }
    // Change between classrooms and update global state
    this.classroomToShow = (this.classroomToShow + 1) % 4
    this.globalState = (this.globalState + 1) % 6

    return classroomMap
  }

  magic = () => {
    if (this.globalState < 2) {
      return this.createCards()
    } else {
      if (this.state.occupation.length !== 0) {
        return this.printClassrooms()
      }
      else {
        this.globalState = 0
        return this.createCards()
      }
    }
  }

  render() {
    return ( // TODO: maybe <tag>{function()}</tag>
      this.magic()
    );
  }

  componentDidMount() {
    this.getReservations()
    this.getOccupation()
    this.timer1 = setInterval(() => {
      this.getReservations()
    }, 10000);
    this.timer2 = setInterval(() => {
      this.getOccupation()
    }, 300000)
  }
  componentWillUnmount() {
    clearInterval(this.timer1);
    clearInterval(this.timer2);
  }
}

export default Main;