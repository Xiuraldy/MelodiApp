<script setup lang="ts">
const users = [
  {
    name: "Juan Pérez",
    role: "Cantante",
    email: "j@gmail.com",
    celphone: "321578962",
  },
  {
    name: "María Gómez",
    role: "Pro. en Vivo",
    email: "m@gmail.com",
    celphone: "373737483",
  },
  {
    name: "Paula Correa",
    role: "Pianista / Baterista",
    email: "p@gmail.com",
    celphone: "378348247",
  },
  {
    name: "Tatiana Rios",
    role: "Líder",
    email: "t@gmail.com",
    celphone: "321508962",
  },
];
</script>

<template>
  <main>
    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th>Nombre</th>
            <th>Rol</th>
            <th>Correo Electrónico</th>
            <th>Celular</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.name">
            <td>{{ user.name }}</td>
            <td>{{ user.role }}</td>
            <td>{{ user.email }}</td>
            <td>{{ user.celphone }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </main>
</template>

<style scoped>
.table-container {
  padding: 20px;
  display: flex;
  justify-content: center;
}

table {
  width: 90%;
  border-collapse: collapse;
  background-color: var(--white);
  border-radius: 12px;
  overflow: hidden;
  margin: 20px auto;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  font-family: 'Arial', sans-serif;
}

thead {
  background-color: var(--color-tertiary);
  color: var(--white);
  text-transform: uppercase;
  font-weight: bold;
}

th, td {
  padding: 15px;
  text-align: left;
}

th {
  letter-spacing: 1px;
}

tbody tr:nth-child(odd) {
  background-color: var(--white);
}

tbody tr:nth-child(even) {
  background-color: #f3f3f3;
}

tbody tr:hover {
  background-color: #e8f5e9;
  transition: background-color 0.3s ease;
}

ul {
  margin: 0;
  padding: 0;
  list-style: none;
}

li {
  padding: 5px 0;
}

li:not(:last-child) {
  border-bottom: 1px dashed #ccc;
}

@media (max-width: 600px) {
  .table-container {
    justify-content: flex-start;
  }
}
</style>