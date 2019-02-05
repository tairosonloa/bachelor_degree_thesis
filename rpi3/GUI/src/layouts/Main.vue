<template>
  <div class="wrapper">
    <header class="header">
        <div class="title">
            <!-- TODO: add logo -->
            Laboratorio de Departamento de Informática
            <Clock />
        </div>
        
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

<script>
import Clock from "~/components/Clock.vue"
export default {
    components: {
        Clock
    }
}
</script>
