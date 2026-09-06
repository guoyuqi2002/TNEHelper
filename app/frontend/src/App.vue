<template>
  <div id="app">
    <MainView />
  </div>
</template>

<script setup>
import MainView from './views/MainView.vue'
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

#app {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  height: 100vh;
  width: 100vw;
  overflow: hidden;
}
</style>
