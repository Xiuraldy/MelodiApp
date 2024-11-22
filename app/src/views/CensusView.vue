<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useCensus } from '@/services/census'
import { useAuthStore } from '@/stores/auth';
import { useConfig } from '@/services/config';
import type { Config } from '@/types';

const { loadCensus, people, totalRecords } = useCensus()
const { loadConfig, updateConfig, userConfig } = useConfig()

let pag = '10';

const selectedAge = ref('');
const selectedWorkclass = ref('');
const selectedFnlwgt = ref('');
const selectedEducation = ref('');
const selectedEducationNum = ref('');
const selectedMaritalStatus = ref('');
const selectedOccupation = ref('');
const selectedRelationship = ref('');
const selectedRace = ref('');
const selectedSex = ref('');
const selectedCapitalGain = ref('');
const selectedCapitalLoss = ref('');
const selectedHoursPerWeek = ref('');
const selectedNativeCountry = ref('');
const selectedIncome = ref('');

const columns = ref([
  { key: 'age', row: 'age' , label: 'Age', options: ['1 to 10', '11 to 20', '21 to 30', '31 to 40', '41 to 50', '51 to 60', '61 to 70', '71 to 80', '81 to 90'], selection: selectedAge, sortOrder: '' },
  { key: 'workclass', row: 'workclass' , label: 'Workclass', options: ['Private', 'Self-emp-not-inc', 'Self-emp-inc', 'Federal-gov', 'Local-gov', 'State-gov', 'Without-pay', 'Never-worked'], selection: selectedWorkclass, sortOrder: '' },
  { key: 'fnlwgt', row: 'fnlwgt' , label: 'Fnlwgt', options: ['1 to 12285', '12286 to 117827', '117828 to 189778', '189779 to 325327', '325328 to 500000'], selection: selectedFnlwgt, sortOrder: '' },
  { key: 'education', row: 'education' , label: 'Education', options: ['Bachelors', 'Some-college', '11th', 'HS-grad', 'Prof-school', 'Assoc-acdm', 'Assoc-voc', '9th', '7th-8th', '12th', 'Masters', '1st-4th', '10th', 'Doctorate', '5th-6th', 'Preschool'], selection: selectedEducation, sortOrder: '' },
  { key: 'educationNum', row: 'education_num' , label: 'Education Num', options: ['1 to 2', '3 to 5', '6 to 8', '9 to 12', '13 to 16'], selection: selectedEducationNum, sortOrder: '' },
  { key: 'maritalStatus', row: 'marital_status' , label: 'Marital Status', options: ['Married-civ-spouse', 'Divorced', 'Never-married', 'Separated', 'Widowed', 'Married-spouse-absent', 'Married-AF-spouse'], selection: selectedMaritalStatus, sortOrder: '' },
  { key: 'occupation', row: 'occupation' , label: 'Occupation', options: ['Tech-support', 'Craft-repair', 'Other-service', 'Sales', 'Exec-managerial', 'Prof-specialty', 'Handlers-cleaners', 'Machine-op-inspct', 'Adm-clerical', 'Farming-fishing', 'Transport-moving', 'Priv-house-serv', 'Protective-serv', 'Armed-Forces'], selection: selectedOccupation, sortOrder: '' },
  { key: 'relationship', row: 'relationship' , label: 'Relationship', options: ['Wife', 'Own-child', 'Husband', 'Not-in-family', 'Other-relative', 'Unmarried'], selection: selectedRelationship, sortOrder: '' },
  { key: 'race', row: 'race' , label: 'Race', options: ['White', 'Asian-Pac-Islander', 'Amer-Indian-Eskimo', 'Other', 'Black'], selection: selectedRace, sortOrder: '' },
  { key: 'sex', row: 'sex' , label: 'Sex', options: ['Male', 'Female'], selection: selectedSex, sortOrder: '' },
  { key: 'capitalGain', row: 'capital_gain' , label: 'Capital Gain', options: ['0 to 0', '1 to 5000', '5001 to 10000', '10001 to 20000'], selection: selectedCapitalGain, sortOrder: '' },
  { key: 'capitalLoss', row: 'capital_loss' , label: 'Capital Loss', options: ['0 to 0', '1 to 100', '101 to 500', '501 to 1000'], selection: selectedCapitalLoss, sortOrder: '' },
  { key: 'hoursPerWeek', row: 'hours_per_week' , label: 'Hours Per Week', options: ['1 to 20', '21 to 40', '41 to 60', '61 to 80'], selection: selectedHoursPerWeek, sortOrder: '' },
  { key: 'nativeCountry', row: 'native_country' , label: 'Native Country', options: ['United-States', 'Cambodia', 'England', 'Puerto-Rico', 'Canada', 'Germany', 'Outlying-US(Guam-USVI-etc)', 'India', 'Japan', 'Greece', 'South', 'China', 'Cuba', 'Iran', 'Honduras', 'Philippines', 'Italy', 'Poland', 'Jamaica', 'Vietnam', 'Mexico', 'Portugal', 'Ireland', 'France', 'Dominican-Republic', 'Laos', 'Ecuador', 'Taiwan', 'Haiti', 'Columbia', 'Hungary', 'Guatemala', 'Nicaragua', 'Scotland', 'Thailand', 'Yugoslavia', 'El-Salvador', 'Trinadad&Tobago', 'Peru', 'Hong', 'Holand-Netherlands'], selection: selectedNativeCountry, sortOrder: '' },
  { key: 'income', row: 'income' , label: 'Income', options: ['<=50K', '>50K'], selection: selectedIncome, sortOrder: '' },
]);

function applyFilters() {
  const filters: Record<string, string> = {};

  let sortBy = ''
  let sortOrder = ''

  columns.value.forEach((column) => {
    if (column.selection) {
      filters[column.key] = column.selection; 
      console.log("filters[column.key]",column.key, filters[column.key])
    }
    if (column.sortOrder) {
      sortBy = column.row; 
      sortOrder = column.sortOrder;
    }
  });

  updateConfig({ ...filters, sortBy, sortOrder, paginator: pag });
  loadCensus({ ...filters, sortBy, sortOrder, paginator: pag })
}

function toggleSortOrder(column: any, order: 'asc' | 'desc') {
  columns.value.forEach((col) => {
    if (col.key !== column.key) {
      col.sortOrder = '';
    }
  });

  column.sortOrder = order;
  pag = '10'
  applyFilters(); 
}

function changePages(direction: string) {
  var pagInt = parseInt(pag)

  if(direction == "back") {
    pagInt = pagInt - 10
    if (pagInt < 10) {
      pagInt = 10
    }
  } else {
    pagInt = pagInt + 10
    if (pagInt > totalRecords.value) {
      pagInt = totalRecords.value
    }
  }

  pag = pagInt.toString()

  applyFilters()
}

onMounted(async () => {
  await loadConfig()

  const filters: Record<string, string> = Object.fromEntries(
    Object.entries(userConfig.value).map(([key, value]) => [key, String(value)])
  );  

  if (filters) {
    columns.value.forEach((column) => {
      if (filters[column.key as keyof Config]) {
        column.selection = String(filters[column.key as keyof Config]);
      }
    });

    pag = filters.paginator || '10';

    const sortBy = columns.value.find((col) => col.row === filters.sortBy);
    if (sortBy) {
      sortBy.sortOrder = filters.sortOrder || '';
    }

    loadCensus({ ...filters })
    return 
  }

   loadCensus({ paginator: pag });
});

console.log("totalRecords",totalRecords.value)

</script>
<template>
  <section>
    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th v-for="column in columns">
              {{ column.label }}
              <select v-model="column.selection" @change="() => { 
                applyFilters()
                pag = '10'
              }">
                <option value="" selected>Filter</option>
                <option v-for="option in column.options.sort((a, b) => parseInt(a.split(' ')[0]) - parseInt(b.split(' ')[0]))" :key="option" :value="option">
                  {{ option }}
                </option>
              </select>
              <div class="content-buttons-order">
                <button class="order-button" @click="toggleSortOrder(column, 'desc')">🡅</button>
                <button class="order-button" @click="toggleSortOrder(column, 'asc')">🡇</button>
              </div>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="person in people">
            <td>{{ person.age }}</td>
            <td>{{ person.workclass }}</td>
            <td>{{ person.fnlwgt }}</td>
            <td>{{ person.education }}</td>
            <td>{{ person['education-num'] }}</td>
            <td>{{ person['marital-status'] }}</td>
            <td>{{ person.occupation }}</td>
            <td>{{ person.relationship }}</td>
            <td>{{ person.race }}</td>
            <td>{{ person.sex }}</td>
            <td>{{ person['capital-gain'] }}</td>
            <td>{{ person['capital-loss'] }}</td>
            <td>{{ person['hours-per-week'] }}</td>
            <td>{{ person['native-country'] }}</td>
            <td>{{ person.income }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="content-buttons-paginator">
      <div class="content-buttons-pag">
        <button class="pag-button" @click="changePages('back')">⬅</button>
        <button class="pag-button" @click="changePages('next')">⮕</button>
      </div>
      <p v-if="totalRecords">Page {{ parseInt(pag)/10 }} of {{ Math.ceil(totalRecords/10) }}</p>
      <p v-else>Records Not Found :c</p>
    </div>
  </section>
</template>

<style>
 .table-container {
    margin-left: 20px;
    overflow-x: auto;
    width: -webkit-fill-available;
  }

  table {
    border-collapse: collapse;
    font-size: 16px;
    text-align: left;
    width: max-content;
  }

  th, td {
    border: 1px solid #ddd;
    padding: 8px;
  }

  th {
    background-color: #f2f2f2;
    color: #333;
    font-weight: bold;
    line-height: 1;
  }

  tr:nth-child(even) {
    background-color: #f9f9f9;
  }

  tr:hover {
    background-color: #f1f1f1;
  }

  th, td {
    text-align: center;
  }
  
  select {
    appearance: none;
    border: 2px solid #e0e0e0;
    border-radius: 5px;
    padding: 5px;
    color: #4a4a4a;
    transition: all 0.3s ease;
    cursor: pointer;
    margin-top: 5px;
    font-family: poppins;
  }

  select:after {
    content: '▼'; 
    font-size: 12px; 
    color: #888; 
  }

  select:hover, select:focus {
    border-color: var(--color-secundary);
    box-shadow: 0 0 8px rgba(176, 176, 255, 0.5);
    outline: none;
  }

  option {
    background-color: #f4f4f9;
    color: #4a4a4a;
  }

  .order-button {
    border: none;
  }

  .content-buttons-order {
    margin: -8px;
    margin-top: 8px;
    background-color: var(--color-secundary);
  }

  .content-buttons-order button:hover {
    background-color: var(--color-secundary);
    color: #fff;
  }

  .content-buttons-order button:active {
    position: relative;
    bottom: 0px;
  }

  .content-buttons-paginator {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .content-buttons-pag {
    display: flex;
    justify-content: center;
  }

  .content-buttons-pag button {
    background-color: var(--color-primary);
    font-size: x-large;
    width: 150px;
  }

  .content-buttons-pag button:hover {
    background-color: var(--color-secundary);
    width: 190px;
  }


  @media (max-width: 600px) {
    th, td {
      padding: 6px;
    }

    .content-buttons-paginator {
      margin-top: 20px;
    }

    .content-buttons-order {
      margin: 6px -6px -6px -6px;
    }
  }
</style>
