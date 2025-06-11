<script setup lang="ts">
</script>

<template>
  <main>
    <div class="search">
      <input type="text" placeholder="Buscar Canción">
    </div>
    <div class="paginator-songs">

      <div class="paginator">< Página 1 de 20 ></div>
      <div class="songs">
        <div class="song">
          <div class="inf-song">
            <div class="name">Danzando</div>
            <div class="author">Christine D’Clario</div>
            <div class="inf">+</div>
          </div>
          <div class="video-song">
            <iframe src="https://www.youtube.com/embed/_1MfaKz_eCQ?si=6RJtcoba1gnNhmhf" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
          </div>
        </div>
        <div class="song">
          <div class="inf-song">
            <div class="name">Danzando</div>
            <div class="author">Christine D’Clario</div>
            <div class="inf">+</div>
          </div>
          <div class="video-song">
            <iframe src="https://www.youtube.com/embed/_1MfaKz_eCQ?si=6RJtcoba1gnNhmhf" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
          </div>
        </div>
        <div class="song">
          <div class="inf-song">
            <div class="name">Danzando</div>
            <div class="author">Christine D’Clario</div>
            <div class="inf">+</div>
          </div>
          <div class="video-song">
            <iframe src="https://www.youtube.com/embed/_1MfaKz_eCQ?si=6RJtcoba1gnNhmhf" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
          </div>
        </div>
        <div class="song">
          <div class="inf-song">
            <div class="name">Danzando</div>
            <div class="author">Christine D’Clario</div>
            <div class="inf">+</div>
          </div>
          <div class="video-song">
            <iframe src="https://www.youtube.com/embed/_1MfaKz_eCQ?si=6RJtcoba1gnNhmhf" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
          </div>
        </div>
      </div>
    </div>
</main>
</template>

<style scoped>

main {
  display: flex;
  flex-direction: column;
  height: 84vh;
  justify-content: space-around;
}

.search {
  display: flex;
  justify-content: center;
}

.search input {
  width: 80%;
  margin: 20px;
  border-radius: 3px;
  height: 30px;
}

.paginator-songs {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-right: 20px;
  margin-left: 20px;
}

.songs {
  display: flex;
  flex-wrap: wrap;
  height: 60vh;
  overflow-y: auto;
  justify-content: center;
}

.song {
  display: flex;
  border: 1px solid;
  margin: 10px;
}

.inf-song {
  margin: 20px;
}

.name {
  font-weight: 700;
}

.author {

}

.video-song {
  display: flex;
}

.inf {
  background-color: var(--color-secundary);
  text-align: center;
  padding: 10px;
  border-radius: 100%;
  font-weight: 600;
  color: var(--white);
  cursor: pointer;
  display: flex;
  justify-content: center;
  width: min-content;
}

@media (max-width: 600px) {
  .video-song {
        width: 10rem;
      }
}
</style>