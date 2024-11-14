<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useEntries } from '@/services/entries'
import type { Entry } from '@/types';

const { loadEntries, deleteEntry, createEntry, editEntry, entries } = useEntries()

const newEntry = reactive({
  title: '',
  content: '',
})

const editMode = ref(false)
const currentEntry = reactive({
  ID: null as number | null,
  title: '',
  content: '',
})

function setEditEntry(entry: Entry) {
  currentEntry.ID = entry.ID
  currentEntry.title = entry.title
  currentEntry.content = entry.content
  editMode.value = true
}

async function saveEditEntry() {
  if (currentEntry.ID) {
    await editEntry(currentEntry.ID, {
      title: currentEntry.title,
      content: currentEntry.content,
    })
    editMode.value = false
    loadEntries()
  }
}

onMounted(loadEntries)
</script>
<template>
  <section>
    <h1>Entries</h1>
    <ul>
      <li v-for="entry in entries" :key="entry.ID">
        {{ entry.ID }}. {{ entry.title }} - {{ entry.content }} 
        <button @click="() => deleteEntry(entry.ID)">Borrar</button>
        <button @click="() => setEditEntry(entry)">Editar</button>
      </li>
    </ul>
    <h2 v-if="!editMode">Add Entry</h2>
    <form v-if="!editMode" @submit.prevent="() => createEntry(newEntry)">
      <input v-model="newEntry.title" type="text" placeholder="Title" />
      <input v-model="newEntry.content" type="text" placeholder="Content" />
      <button type="submit">Add Entry</button>
    </form>

    <h2 v-else>Edit Entry</h2>
    <form v-if="editMode" @submit.prevent="saveEditEntry">
      <input v-model="currentEntry.title" type="text" placeholder="Title" />
      <input v-model="currentEntry.content" type="text" placeholder="Content" />
      <button type="submit">Save Changes</button>
      <button type="button" @click="editMode = false">Cancel</button>
    </form>
  </section>
</template>

<style>
@media (min-width: 1024px) {
  .about {
    min-height: 100vh;
    display: flex;
    align-items: center;
  }
}
</style>
