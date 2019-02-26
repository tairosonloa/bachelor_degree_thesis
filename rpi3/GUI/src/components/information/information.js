import React from "react"
import styles from "./information.module.css"

class Information extends React.Component {
  constructor(props) {
    super(props);
    this.message = "Prueba"
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
    return (
      <div className={styles.marquee}>
        <p>La temperatura en el CPD es de {this.state.cpdStatus["temperature"]} ºC</p>
        <p>La humedad en el CPD está al {this.state.cpdStatus["humidity"]} %</p>
        <p>El estado de la batería del SAI es {this.state.cpdStatus["ups status (LDI rack)"]}.</p>
      </div>
    );
  }
}

export default Information;