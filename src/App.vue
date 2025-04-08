<template>
  <div class="app">
    <div class="input-section">
      <textarea 
        v-model="query" 
        placeholder="请输入查询内容..."
        @keyup.enter="sendQuery"
      ></textarea>
      <button @click="sendQuery">发送</button>
    </div>
    <div class="graph-container" ref="graphContainer"></div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import * as d3 from 'd3'
import axios from 'axios'

const query = ref('')
const graphContainer = ref(null)
const graphData = ref(null)

const sendQuery = async () => {
  try {
    const response = await axios.post('http://192.168.5.6/v1/chat-messages', {
      inputs: {},
      query: query.value,
      response_mode: 'streaming',
      conversation_id: '',
      user: 'abc-123'
    }, {
      headers: {
        'Authorization': 'Bearer app-aQmrylCiemu3Ksg2xrLexgG5',
        'Content-Type': 'application/json'
      }
    })

    // 处理流式响应
    const reader = response.data.getReader()
    let result = ''
    
    while (true) {
      const { done, value } = await reader.read()
      if (done) break
      result += new TextDecoder().decode(value)
    }

    // 解析JSON并更新图形
    try {
      const jsonData = JSON.parse(result)
      graphData.value = jsonData
      drawGraph()
    } catch (e) {
      console.error('解析JSON失败:', e)
    }
  } catch (error) {
    console.error('请求失败:', error)
  }
}

const drawGraph = () => {
  if (!graphData.value || !graphContainer.value) return

  // 清除现有图形
  d3.select(graphContainer.value).selectAll('*').remove()

  // 创建SVG容器
  const width = graphContainer.value.clientWidth
  const height = graphContainer.value.clientHeight
  const svg = d3.select(graphContainer.value)
    .append('svg')
    .attr('width', width)
    .attr('height', height)

  // 创建力导向图
  const simulation = d3.forceSimulation(graphData.value.nodes)
    .force('link', d3.forceLink(graphData.value.links).id(d => d.id))
    .force('charge', d3.forceManyBody())
    .force('center', d3.forceCenter(width / 2, height / 2))

  // 绘制连线
  const links = svg.append('g')
    .selectAll('line')
    .data(graphData.value.links)
    .enter()
    .append('line')
    .attr('stroke', '#999')
    .attr('stroke-width', 2)

  // 绘制节点
  const nodes = svg.append('g')
    .selectAll('circle')
    .data(graphData.value.nodes)
    .enter()
    .append('circle')
    .attr('r', 10)
    .attr('fill', '#69b3a2')

  // 添加节点标签
  const labels = svg.append('g')
    .selectAll('text')
    .data(graphData.value.nodes)
    .enter()
    .append('text')
    .text(d => d.id)
    .attr('font-size', 12)
    .attr('dx', 15)
    .attr('dy', 4)

  // 更新位置
  simulation.on('tick', () => {
    links
      .attr('x1', d => d.source.x)
      .attr('y1', d => d.source.y)
      .attr('x2', d => d.target.x)
      .attr('y2', d => d.target.y)

    nodes
      .attr('cx', d => d.x)
      .attr('cy', d => d.y)

    labels
      .attr('x', d => d.x)
      .attr('y', d => d.y)
  })
}

onMounted(() => {
  // 监听窗口大小变化
  window.addEventListener('resize', drawGraph)
})
</script>

<style scoped>
.app {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.input-section {
  padding: 20px;
  display: flex;
  gap: 10px;
}

textarea {
  flex: 1;
  height: 100px;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
  resize: none;
}

button {
  padding: 10px 20px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:hover {
  background-color: #45a049;
}

.graph-container {
  flex: 1;
  background-color: #f5f5f5;
}
</style> 