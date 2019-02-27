import React from "react"

import Clock from "../components/clock/clock.js"
import Cards from "../components/cards/cards.js"
import Information from "../components/information/information.js"

export default () => (
  <div className="wrapper">
    <header className="header">
      <div className="headerTitle">
        {/* TODO: add logo */}
        <h1>Laboratorio de Departamento de Informática</h1>
      </div>
      <div className="datetime">
        <Clock/>
      </div>
    </header>
    <main className="main">

    </main>
    <aside className="aside">
      <Cards/>
    </aside>
    <footer className="footer">
      <Information />
    </footer>
  </div>
)
