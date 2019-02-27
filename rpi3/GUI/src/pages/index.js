import React from "react"
import { Helmet } from "react-helmet"

import Clock from "../components/clock/clock.js"
import Cards from "../components/cards/cards.js"
import Information from "../components/information/information.js"

export default () => (
  <div className="wrapper">
    <Helmet>
      <meta charSet="utf-8" />
      <title>LDI</title>
      <meta http-equiv="Content-Language" content="es" />
    </Helmet>
    <header className="header">
      <div className="headerTitle">
        {/* TODO: add logo */}
        <h1>Laboratorio del Departamento de Informática</h1>
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
