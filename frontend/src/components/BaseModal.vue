<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content">
      <header class="modal-header">
        <slot name="header"></slot>
      </header>
      <section class="modal-body">
        <slot name="body"></slot>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
// Простая реализация, можно расширить по необходимости.
</script>

<style>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: fit-content;
}

.modal-content {
  background: #d9d9d9;
  padding: 20px;
  border-radius: 5px;
  max-width: 400px;
  width: 100%;
  font-size: large;
}

.modal-header h3 {
  padding: 0;
}

.modal-body {
  max-height: 400px;
  overflow-y: auto;
}

.modal-group {
  display: flex;
  flex-direction: column;
  max-width: 400px;
  width: 95%;
  margin: 0 auto;
}

.modal-group p {
  padding: 0 0;
  margin: 5px;
}

.modal-group-button {
  margin-top: 10px;
  justify-self: end;
}

.modal-group label {
  font-weight: 500;
  margin-bottom: 0.25rem;
  align-self: start;
}

.form-group input,
.form-group select,
.form-group textarea {
  padding: 0.5rem;
  font-size: 1rem;
  border: 1px solid #ccc;
  border-radius: 4px;
}
</style>
