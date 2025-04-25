<template>
  <nav>
    <div>
      <router-link to="/runs">Runs</router-link>
    </div>

    <div>v0.2.3</div>
  </nav>
  <router-view />
</template>

<style>
html,
body,
#app {
  font-family: "Poppins", sans-serif;
  color: #40231d;
  background-color: #d9d9d9;
}

nav {
  padding: 20px 50px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

nav a {
  font-size: large;
  font-weight: bold;
  color: #40231d;
  margin-right: 10px;
}

nav a.router-link-exact-active {
  color: #a6523f;
}

.material-symbols-outlined {
  font-variation-settings: "FILL" 0, "wght" 400, "GRAD" 0, "opsz" 24;
}

.view-header {
  display: flex;
  justify-content: space-between;
}

.view-header h1 {
  justify-self: start;
  padding: 0 50px;
  margin: 0;
}

.view-header button {
  padding: 0 60px;
  margin-top: 5px;
}

.view-header .material-symbols-outlined {
  font-size: 40px;
}

button {
  background-color: #d9d9d9;

  font-size: 20px;

  border: none;

  padding: 1px 1px;

  color: #40231d;
}
input {
  max-width: 100%;
  font-size: medium;
}
select {
  max-width: 100%;
  font-size: medium;
}
textarea {
  max-width: 100%;
  resize: none;
}
button:hover {
  color: #a6523f;
}
h3 {
  display: flex;
  justify-self: start;
  padding: 0 20px;
}
.sendpost-list {
  margin-top: 10px;
  padding: 0 50px;
  min-width: fit-content;
}
.sendpost-header {
  display: flex;
  align-items: center;
}
.sendpost-header span {
  margin-right: 10px;
}
.sendpost-header strong {
  color: #a6523f;
}
.actions {
  margin-left: auto;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
}

.actions button {
  vertical-align: middle;
}
</style>
