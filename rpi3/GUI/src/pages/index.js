import React from "react"

import Clock from "../components/clock.js"

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
    <article className="article">

    </article>
    <aside className="aside">
    
    </aside>
    <footer className="footer">

    </footer>
  </div>
)
