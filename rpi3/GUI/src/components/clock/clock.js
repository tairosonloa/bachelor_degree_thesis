import React from "react"
import styles from "./clock.module.css"

class Clock extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      time : this.getCurrentTime(),
      date : this.getCurrentDate()
    }
  }
  /**
   * Gets current date
   */
  getCurrentDate = () => {
    return new Date().toLocaleDateString("es-ES", { weekday: 'long', month: 'long', day: 'numeric' });
  }

  /**
   * Gets current time
   */
  getCurrentTime = () => {
    let time = new Date().toLocaleTimeString("es-ES", {hour: "2-digit", minute:"2-digit"});
    if (time === "00:00") {
      this.setState({date : this.getCurrentDate()})
    }
    return time;
  }

  componentDidMount() {
    this.timer = setInterval(() => {
      this.setState({time : this.getCurrentTime()});
    }, 10000);
  }

  componentWillUnmount() {
    clearInterval(this.timer);
  }
  
  render() {
    return (
      <div className={styles.clock}>
        <div className={styles.date}>{this.state.date},</div>
        <div className={styles.time}>{this.state.time}</div>
      </div>
    );
  }
}

export default Clock;