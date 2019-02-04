<template>
  <div class="wrapper">
    <header class="header">
      Laboratorio de Departamento de Informática
    </header>
    <article class="article">
        <slot />
    </article>
    <aside class="aside">
        <slot name="reservations" />
    </aside>
    <footer class="footer">
        <slot name="titulars" />
    </footer>
  </div>
</template>