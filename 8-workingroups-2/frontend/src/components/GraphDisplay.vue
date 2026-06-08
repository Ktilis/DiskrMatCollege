<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { Network } from 'vis-network'
import { DataSet } from 'vis-data'
import type { Node, Edge, Options } from 'vis-network'

const props = defineProps<{
  graphData: Record<string, { Node: string; Weight: number }[]>
  activePaths: string[][]
  isSearching: boolean
}>()

const emit = defineEmits(['select-nodes', 'animation-finished'])

const container = ref<HTMLElement | null>(null)
let network: Network | null = null
const nodes = new DataSet<Node>([])
const edges = new DataSet<Edge>([])

const selectedStart = ref<string | null>(null)
const selectedEnd = ref<string | null>(null)
const isAborted = ref(false)

const initGraph = () => {
  if (!container.value) return

  nodes.clear()
  edges.clear()

  const nodeSet = new Set<string>()
  const edgeList: { from: string; to: string }[] = []

  Object.entries(props.graphData).forEach(([from, connections]) => {
    nodeSet.add(from)
    connections.forEach(conn => {
      nodeSet.add(conn.Node)
      // Добавляем ребро только один раз для неориентированного графа
      if (from < conn.Node) {
        edgeList.push({ from, to: conn.Node })
      }
    })
  })

  nodes.add(Array.from(nodeSet).map(id => ({
    id,
    label: id,
    color: {
      background: '#1e1e1e',
      border: '#7c4dff',
      highlight: { background: '#7c4dff', border: '#b47cff' }
    },
    font: { color: '#ffffff' }
  })))

  edges.add(edgeList.map(e => ({
    from: e.from,
    to: e.to,
    color: { color: '#3f1dcb', highlight: '#7c4dff' },
    width: 2
  })))

  const data = { nodes, edges }
  const options: Options = {
    physics: {
      enabled: true,
      barnesHut: {
        gravitationalConstant: -2000,
        centralGravity: 0.3,
        springLength: 95,
        springConstant: 0.04,
        damping: 0.09,
        avoidOverlap: 0
      },
      stabilization: { iterations: 150 }
    },
    interaction: {
      hover: true,
      selectConnectedEdges: false
    }
  }

  network = new Network(container.value, data, options)

  network.on('click', (params) => {
    if (params.nodes.length > 0) {
      const nodeId = params.nodes[0] as string
      if (!selectedStart.value) {
        selectedStart.value = nodeId
      } else if (!selectedEnd.value || selectedEnd.value === nodeId) {
        selectedEnd.value = nodeId
      } else {
        selectedStart.value = nodeId
        selectedEnd.value = null
      }
      emit('select-nodes', { start: selectedStart.value, end: selectedEnd.value })
    }
  })
}

const highlightPath = async (path: string[]) => {
  if (isAborted.value) return
  resetHighlights()
  
  for (let i = 0; i < path.length; i++) {
    if (isAborted.value) break
    const nodeId = path[i]
    nodes.update({
      id: nodeId,
      color: { background: '#ff5252', border: '#ff8a80' },
      font: { color: '#ffffff' }
    })

    if (i > 0) {
      const from = path[i-1]
      const to = path[i]
      const edgeId = edges.get().find(e => 
        (e.from === from && e.to === to) || (e.from === to && e.to === from)
      )?.id
      
      if (edgeId) {
        edges.update({
          id: edgeId,
          color: { color: '#ff5252' },
          width: 4
        })
      }
    }
    await new Promise(resolve => setTimeout(resolve, 400))
  }
}

const resetHighlights = () => {
  nodes.get().forEach(node => {
    const isStart = node.id === selectedStart.value
    const isEnd = node.id === selectedEnd.value
    nodes.update({
      id: node.id,
      color: {
        background: isStart || isEnd ? '#7c4dff' : '#1e1e1e',
        border: '#7c4dff'
      }
    })
  })
  edges.get().forEach(edge => {
    edges.update({
      id: edge.id,
      color: { color: '#3f1dcb' },
      width: 2
    })
  })
}

watch(() => props.graphData, initGraph)

watch(() => props.activePaths, async (newPaths) => {
  if (newPaths.length === 0) {
    resetHighlights()
    return
  }
  
  isAborted.value = false
  for (const path of newPaths) {
    if (isAborted.value) break
    await highlightPath(path)
    if (newPaths.length > 1 && !isAborted.value) {
      await new Promise(resolve => setTimeout(resolve, 1000))
    }
  }
  emit('animation-finished')
})

watch([selectedStart, selectedEnd], () => {
  resetHighlights()
})

onMounted(initGraph)
onUnmounted(() => {
  if (network) network.destroy()
})

defineExpose({
  resetSelection: () => {
    selectedStart.value = null
    selectedEnd.value = null
    resetHighlights()
  },
  abortAnimation: () => {
    isAborted.value = true
    resetHighlights()
  }
})
</script>

<template>
  <div class="graph-container" ref="container"></div>
</template>

<style scoped>
.graph-container {
  width: 100%;
  height: 600px;
  background: var(--surface);
  border-radius: 12px;
  box-shadow: inset 0 0 20px rgba(0,0,0,0.5);
}
</style>
