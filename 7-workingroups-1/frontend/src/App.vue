<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import TreeDisplay from './components/TreeDisplay.vue'

interface TreeNode {
  Val: number
  Left: TreeNode | null
  Right: TreeNode | null
}

interface TreeResponse {
  sessionId: string
  tree: TreeNode
}

interface SearchResponse {
  path: number[]
  found: boolean
}

const tree = ref<TreeNode | null>(null)
const sessionId = ref('')
const targetValue = ref(0)
const activeNodeVal = ref<number | null>(null)
const finalNodeVal = ref<number | null>(null)
const isSearching = ref(false)
const searchFinished = ref(false)

const fetchTree = async () => {
  try {
    const response = await axios.get<TreeResponse>('/api/tree')
    tree.value = response.data.tree
    sessionId.value = response.data.sessionId
    activeNodeVal.value = null
    finalNodeVal.value = null
    searchFinished.value = false
  } catch (error) {
    console.error('Ошибка при получении дерева:', error)
  }
}

const startSearch = async () => {
  if (!sessionId.value || isSearching.value) return
  
  isSearching.value = true
  searchFinished.value = false
  activeNodeVal.value = null
  finalNodeVal.value = null

  try {
    const response = await axios.post<SearchResponse>('/api/search', {
      sessionId: sessionId.value,
      target: targetValue.value
    })
    
    const path = response.data.path
    const found = response.data.found

    for (let i = 0; i < path.length; i++) {
      activeNodeVal.value = path[i]
      await new Promise(resolve => setTimeout(resolve, 600))
    }

    if (found) {
      finalNodeVal.value = path[path.length - 1]
    } else {
      activeNodeVal.value = null
    }
    
    searchFinished.value = true
  } catch (error) {
    console.error('Ошибка во время поиска:', error)
  } finally {
    isSearching.value = false
  }
}

onMounted(fetchTree)
</script>

<template>
  <div class="app-container">
    <header>
      <h1>Визуализация поиска</h1>
      <div class="controls">
        <div class="input-group">
          <label>Искомое значение:</label>
          <input type="number" v-model="targetValue" :disabled="isSearching" />
        </div>
        <button @click="startSearch" :disabled="isSearching" class="btn-search">Поиск</button>
        <button @click="fetchTree" :disabled="isSearching" class="btn-reset">Новое дерево</button>
      </div>
    </header>

    <main>
      <div v-if="tree" class="tree-canvas">
        <TreeDisplay 
          :node="tree" 
          :active-val="activeNodeVal" 
          :final-val="finalNodeVal"
          :is-finished="searchFinished"
        />
      </div>
      <div v-else class="loading">Генерация дерева...</div>
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
  margin-bottom: 3rem;
  width: 100%;
}

h1 {
  color: var(--primary-light);
  margin-bottom: 2rem;
  font-weight: 300;
  letter-spacing: 2px;
}

.controls {
  display: inline-flex;
  gap: 1.5rem;
  align-items: center;
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

input {
  background: var(--bg);
  border: 1px solid var(--primary-dark);
  color: white;
  padding: 0.6rem;
  border-radius: 4px;
  width: 80px;
  outline: none;
}

input:focus {
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

.tree-canvas {
  width: 100%;
  flex-grow: 1;
  display: flex;
  justify-content: center;
  padding: 2rem 0;
  overflow-x: auto;
}

.loading {
  color: var(--primary-light);
  font-style: italic;
}
</style>
