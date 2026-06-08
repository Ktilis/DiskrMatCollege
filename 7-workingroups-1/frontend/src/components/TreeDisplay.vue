<script setup lang="ts">
import { computed } from 'vue'

interface TreeNode {
  Val: number
  Left: TreeNode | null
  Right: TreeNode | null
}

interface LayoutNode {
  node: TreeNode
  x: number
  y: number
  parentX: number | null
  parentY: number | null
}

const props = defineProps<{
  node: TreeNode | null
  activeVal: number | null
  finalVal: number | null
  isFinished: boolean
}>()

const NODE_R = 20
const LEVEL_H = 70
const GAP = 14

function subtreeWidth(node: TreeNode | null): number {
  if (!node) return NODE_R * 2
  if (!node.Left && !node.Right) return NODE_R * 2
  const lw = node.Left  ? subtreeWidth(node.Left)  : 0
  const rw = node.Right ? subtreeWidth(node.Right) : 0
  if (node.Left && node.Right) return lw + GAP + rw
  return Math.max(NODE_R * 2, (node.Left ? lw : rw) + NODE_R * 2)
}

function collectNodes(
  node: TreeNode | null, 
  x: number, 
  y: number, 
  parentX: number | null, 
  parentY: number | null, 
  list: LayoutNode[] = []
): LayoutNode[] {
  if (!node) return list
  list.push({ node, x, y, parentX, parentY })

  const lw = subtreeWidth(node.Left)
  const rw = subtreeWidth(node.Right)

  if (node.Left && node.Right) {
    collectNodes(node.Left,  x - rw / 2 - GAP / 2, y + LEVEL_H, x, y, list)
    collectNodes(node.Right, x + lw / 2 + GAP / 2, y + LEVEL_H, x, y, list)
  } else if (node.Left) {
    collectNodes(node.Left,  x - lw / 2, y + LEVEL_H, x, y, list)
  } else if (node.Right) {
    collectNodes(node.Right, x + rw / 2, y + LEVEL_H, x, y, list)
  }
  return list
}

const layout = computed(() => {
  if (!props.node) return { nodes: [], width: 0, height: 0 }

  const totalWidth = subtreeWidth(props.node)
  const startX = totalWidth / 2 + NODE_R + 10
  const nodes = collectNodes(props.node, startX, NODE_R + 10, null, null)

  const PADDING = NODE_R + 10
  const minX = Math.min(...nodes.map(n => n.x))
  const maxX = Math.max(...nodes.map(n => n.x))
  const maxY = Math.max(...nodes.map(n => n.y))

  // Сдвиг всех узлов, чтобы левый край крайнего левого круга начинался с PADDING
  const offsetX = PADDING - minX + NODE_R
  nodes.forEach(n => { 
    n.x += offsetX; 
    if (n.parentX !== null) n.parentX += offsetX 
  })

  const width  = maxX - minX + NODE_R * 2 + PADDING * 2
  const height = maxY + NODE_R + PADDING

  return { nodes, width, height }
})
</script>

<template>
  <svg
      v-if="node"
      :width="layout.width"
      :height="layout.height"
      :viewBox="`0 0 ${layout.width} ${layout.height}`"
      class="tree-svg"
  >
    <!-- Сначала ребра (под узлами) -->
    <g class="edges">
      <line
          v-for="(item, i) in layout.nodes.filter(n => n.parentX !== null)"
          :key="'e' + i"
          :x1="item.parentX ?? 0"
          :y1="item.parentY ?? 0"
          :x2="item.x"
          :y2="item.y"
          class="tree-edge"
      />
    </g>

    <!-- Узлы -->
    <g
        v-for="(item, i) in layout.nodes"
        :key="'n' + i"
        :class="[
        'tree-node',
        { 'is-active':  activeVal === item.node.Val },
        { 'is-final':   isFinished && finalVal === item.node.Val }
      ]"
    >
      <circle
          :cx="item.x"
          :cy="item.y"
          :r="NODE_R"
      />
      <text
          :x="item.x"
          :y="item.y"
          text-anchor="middle"
          dominant-baseline="central"
          class="node-label"
      >{{ item.node.Val }}</text>
    </g>
  </svg>
</template>

<style scoped>
.tree-svg {
  overflow: visible;
  display: block;
}

.tree-edge {
  stroke: var(--primary-dark, #534AB7);
  stroke-width: 1.5px;
  fill: none;
  transition: stroke 0.3s ease;
}

.tree-node circle {
  fill: var(--surface, #1a1a2e);
  stroke: var(--primary, #7C4DFF);
  stroke-width: 2px;
  transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  filter: drop-shadow(0 0 6px rgba(124, 77, 255, 0.2));
}

.node-label {
  font-size: 13px;
  font-weight: bold;
  fill: var(--text, #e0e0e0);
  pointer-events: none;
  transition: fill 0.3s ease;
}

/* Активное состояние */
.tree-node.is-active circle {
  fill: var(--primary, #7C4DFF);
  stroke: var(--primary-light, #b39dff);
  filter: drop-shadow(0 0 12px rgba(124, 77, 255, 0.8));
}

.tree-node.is-active .node-label {
  fill: #fff;
}

/* Конечное состояние */
.tree-node.is-final circle {
  fill: var(--accent, #ff5252);
  stroke: #ff8a80;
  filter: drop-shadow(0 0 16px rgba(255, 82, 82, 0.8));
}

.tree-node.is-final .node-label {
  fill: #fff;
}
</style>
