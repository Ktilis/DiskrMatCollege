<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import GraphDisplay from './components/GraphDisplay.vue'

interface GraphResponse {
  sessionId: string
  graph: Record<string, { Node: string; Weight: number }[]>
}

interface SearchResponse {
  paths: string[][]
}

const graphData = ref<Record<string, { Node: string; Weight: number }[]>>({})
const sessionId = ref('')
const startUser = ref('')
const endUser = ref('')
const algorithm = ref('bfs')
const activePaths = ref<string[][]>([])
const isSearching = ref(false)
const isAnimating = ref(false)
const graphDisplayRef = ref<any>(null)

const fetchGraph = async () => {
  try {
    const response = await axios.get<GraphResponse>('/api/graph')
    graphData.value = response.data.graph
    sessionId.value = response.data.sessionId
    activePaths.value = []
    startUser.value = ''
    endUser.value = ''
    if (graphDisplayRef.value) graphDisplayRef.value.resetSelection()
  } catch (error) {
    console.error('Ошибка при получении графа:', error)
  }
}

const handleSelectNodes = (nodes: { start: string; end: string }) => {
  startUser.value = nodes.start || ''
  endUser.value = nodes.end || ''
}

const startSearch = async () => {
  if (!sessionId.value || isSearching.value || !startUser.value || !endUser.value) return
  
  isSearching.value = true
  activePaths.value = []

  try {
    const response = await axios.post<SearchResponse>('/api/search', {
      sessionId: sessionId.value,
      start: startUser.value,
      end: endUser.value,
      algorithm: algorithm.value
    })
    
    activePaths.value = response.data.paths
    if (activePaths.value.length === 0 || (activePaths.value.length === 1 && activePaths.value[0] === null)) {
      alert('Пути не найдены!')
    } else {
      isAnimating.value = true
    }
  } catch (error) {
    console.error('Ошибка при поиске:', error)
  } finally {
    isSearching.value = false
  }
}

const abortAnimation = () => {
  if (graphDisplayRef.value) {
    graphDisplayRef.value.abortAnimation()
  }
  isAnimating.value = false
}

const onAnimationFinished = () => {
  isAnimating.value = false
}

onMounted(fetchGraph)
</script>

<template>
  <div class="app-container">
    <header>
      <div class="controls">
        <div class="input-group">
          <label>Начало:</label>
          <input type="text" v-model="startUser" :disabled="isSearching || isAnimating" placeholder="Клик на узел" />
        </div>
        <div class="input-group">
          <label>Конец:</label>
          <input type="text" v-model="endUser" :disabled="isSearching || isAnimating" placeholder="Клик на узел" />
        </div>
        <div class="input-group">
          <label>Алгоритм:</label>
          <select v-model="algorithm" :disabled="isSearching || isAnimating">
            <option value="bfs">BFS (Кратчайший путь)</option>
            <option value="dfs">DFS (Все пути)</option>
          </select>
        </div>
        <button v-if="!isAnimating" @click="startSearch" :disabled="isSearching || !startUser || !endUser" class="btn-search">Найти связи</button>
        <button v-else @click="abortAnimation" class="btn-abort">Прервать анимацию</button>
        <button @click="fetchGraph" :disabled="isSearching || isAnimating" class="btn-reset">Новый граф</button>
      </div>
    </header>

    <main>
      <div class="graph-canvas">
        <GraphDisplay 
          ref="graphDisplayRef"
          :graph-data="graphData" 
          :active-paths="activePaths"
          :is-searching="isSearching"
          @select-nodes="handleSelectNodes"
          @animation-finished="onAnimationFinished"
        />
      </div>
      <div v-if="Object.keys(graphData).length === 0" class="loading">Генерация социального графа...</div>
      <div v-if="activePaths.length > 0" class="results-info">
        Найдено путей: {{ activePaths.length }}
      </div>
    </main>
  </div>
</template>

<style>
:root {
  --primary: #7c4dff;
  --primary-light: #b47cff;
  --primary-dark: #3f1dcb;
  --bg: #121212;
  --surface: #1e1e1e;
  --text: #ffffff;
  --accent: #ff5252;
}

body {
  margin: 0;
  background-color: var(--bg);
  color: var(--text);
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.app-container {
  width: 100%;
  margin: 0;
  padding: 2rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  min-height: 100vh;
  box-sizing: border-box;
}

header {
  text-align: center;
  margin-bottom: 2rem;
  width: 100%;
}

h1 {
  color: var(--primary-light);
  margin-bottom: 1.5rem;
  font-weight: 300;
  letter-spacing: 2px;
}

.controls {
  display: inline-flex;
  flex-wrap: wrap;
  gap: 1.5rem;
  align-items: flex-end;
  justify-content: center;
  background: var(--surface);
  padding: 1.5rem 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.3);
}

.input-group {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.5rem;
}

.input-group label {
  font-size: 0.8rem;
  color: var(--primary-light);
  text-transform: uppercase;
}

input, select {
  background: var(--bg);
  border: 1px solid var(--primary-dark);
  color: white;
  padding: 0.6rem;
  border-radius: 4px;
  width: 120px;
  outline: none;
}

select {
  width: 180px;
}

input:focus, select:focus {
  border-color: var(--primary);
}

button {
  padding: 0.8rem 1.5rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: bold;
  transition: all 0.3s;
}

.btn-search {
  background: var(--primary);
  color: white;
}

.btn-search:hover:not(:disabled) {
  background: var(--primary-light);
  transform: translateY(-2px);
}

.btn-abort {
  background: var(--accent);
  color: white;
}

.btn-abort:hover {
  background: #ff8a80;
  transform: translateY(-2px);
}

.btn-reset {
  background: transparent;
  border: 1px solid var(--primary);
  color: var(--primary-light);
}

.btn-reset:hover:not(:disabled) {
  background: rgba(124, 77, 255, 0.1);
}

button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.graph-canvas {
  width: 1000px;
  max-width: 90vw;
  margin-top: 1rem;
}

.loading {
  color: var(--primary-light);
  font-style: italic;
  margin-top: 2rem;
}

.results-info {
  margin-top: 1rem;
  color: var(--primary-light);
  font-weight: bold;
}
</style>
