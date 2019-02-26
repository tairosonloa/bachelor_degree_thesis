import React from "react"

import Clock from "../components/clock.js"
import Reservations from "../components/reservations.js"
import Information from "../components/information/information.js"

export default () => (
  <div className="wrapper">
    <header className="header">
      <div className="headerTitle">
        {/* TODO: add logo */}
        Laboratorio de Departamento de Informática
      </div>
      <div className="datetime">
        <Clock/>
      </div>
    </header>
    <main className="main">

    </main>
    <aside className="aside">
      <Reservations/>
    </aside>
    <footer className="footer">
      <Information />
    </footer>
  </div>
)
