<script setup lang="ts">
</script>

<template>
  <main>
    <div class="programations">
      <div class="card-programation">
        <div class="header-card">
          <div class="date">Diciembre 24</div>
          <div class="role">Cantante</div>
          <div class="event">Evento: Navidad</div>
        </div>
        <hr>
        <div class="songs">
          <div class="title">Canciones:</div>
          <div class="song">Espíritu Santo</div>
          <div class="song">Fiesta en el Desierto</div>
          <div class="song">La Bondad de Dios</div>
        </div>
      </div>
      <div class="card-programation">
        <div class="header-card">
          <div class="date">Enero 30</div>
          <div class="role">Pro. en Vivo</div>
          <div class="event">Evento: Culto General</div>
        </div>
        <hr>
        <div class="songs">
          <div class="title">Canciones:</div>
          <div class="song">Danzando</div>
          <div class="song">Vida Tu Me Das</div>
          <div class="song">Yaweh</div>
        </div>
      </div>
    </div>
  </main>
</template>

<style scoped>
  .programations {
    display: flex;
    justify-content: flex-start;
    flex-wrap: wrap;
    padding: 50px;
  }

  .card-programation {
    border: 2px solid var(--black);
    padding: 20px;
    border-radius: 20px;
    margin: 20px;
    box-shadow: 2px 2px 2px 1px rgba(0, 0, 0, 0.2);
    cursor: pointer;
    transition: .2s all;
  }

  .card-programation:hover {
    border: 2px solid var(--color-secundary);
    position: relative;
    bottom: 2px
  }

  .header-card {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .date {
    font-weight: 700;
    color: var(--color-secundary);
  }

  .event {
    font-weight: 700;
    color: var(--color-tertiary);
  }

  hr {
    margin-top: 10px;
    margin-bottom: 10px;
    border: 1px dashed;
    color: gray;
  }

  .title {
    font-weight: 600;
  }

  .songs {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .song {
    background-color: var(--color-tertiary);
    width: -webkit-fill-available;
    margin-bottom: 5px;
    border-radius: 20px;
    text-align: center;
    padding-right: 20px;
    padding-left: 20px;
    box-shadow: 2px 2px 2px 1px rgba(0, 0, 0, 0.2);
  }
</style>