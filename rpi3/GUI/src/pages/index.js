import React from "react"
import { Helmet } from "react-helmet"

import Clock from "../components/clock/clock.js"
import Footer from "../components/footer/footer.js"
import Classrooms from "../components/classrooms/classrooms.js";

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
    <Classrooms/>
    <footer className="footer">
      <Footer/>
    </footer>
  </div>
)
