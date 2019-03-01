import React from "react"
import styles from "./information.module.css"

class Information extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      cpdStatus : []
    }
    try { // Load config
      this.config = require("/etc/rpi3_conf.json")
    } catch {
      this.config = require("../../../../../config.json")
    }
  }

  getValues = () => {
    fetch("http://" + this.config.Rpi2APIAddress + ":" + this.config.Rpi2APIPort + "/cpd-status")
      .then(response => response.json())
      .then(json => this.setState({cpdStatus: json}))
      .catch(error => console.log('Request failed', error))
  }

  componentDidMount() {
    this.getValues()
    this.timer = setInterval(() => {
      this.getValues()
    }, 60000);
  }
  
  componentWillUnmount() {
    clearInterval(this.timer);
  }
  
  render() {
    let temp = "La temperatura en el CPD es de " + this.state.cpdStatus["temperature"] + " ºC"
    let hum = "La humedad en el CPD está al " + this.state.cpdStatus["humidity"] + " %"
    let sai = "El estado de la batería del SAI es " + this.state.cpdStatus["ups status (LDI rack)"] + "."
    let message = temp + "   •   " + hum + "   •   " + sai + "   •   "
    return (
      <div className={styles.marquee}>
        <pre>{message}</pre>
        <pre>{message}</pre>
        <pre>{message}</pre>
      </div>
    );
  }
}

export default Information;