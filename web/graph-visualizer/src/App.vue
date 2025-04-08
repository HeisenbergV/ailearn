<script setup>
import { ref, onMounted } from 'vue'
import { Input, Button, message, Card } from 'ant-design-vue'
import * as d3 from 'd3'

const inputValue = ref('')
const graphData = ref(null)
const svgRef = ref(null)
const loading = ref(false)
const answerText = ref('')
const streamLogs = ref([])
const showStreamLogs = ref(false)

const addLog = (type, data) => {
  streamLogs.value.push({
    type,
    data,
    timestamp: new Date().toLocaleTimeString()
  })
  // 保持日志在可视区域内
  setTimeout(() => {
    const logContainer = document.querySelector('.stream-container')
    if (logContainer) {
      logContainer.scrollTop = logContainer.scrollHeight
    }
  }, 0)
}

const drawGraph = (data) => {
  if (!svgRef.value) return
  
  // 清除之前的图形
  d3.select(svgRef.value).selectAll('*').remove()
  
  const width = svgRef.value.clientWidth
  const height = svgRef.value.clientHeight
  
  // 创建力导向图
  const simulation = d3.forceSimulation(data.nodes)
    .force('link', d3.forceLink(data.links).id(d => d.id).distance(100))
    .force('charge', d3.forceManyBody().strength(-400))
    .force('center', d3.forceCenter(width / 2, height / 2))
    .force('collision', d3.forceCollide().radius(50))
  
  // 创建 SVG 元素
  const svg = d3.select(svgRef.value)
    .attr('width', width)
    .attr('height', height)
    .style('background-color', '#fafafa')
  
  // 创建箭头
  svg.append('defs').selectAll('marker')
    .data(['end'])
    .enter().append('marker')
    .attr('id', d => d)
    .attr('viewBox', '0 -5 10 10')
    .attr('refX', 20)
    .attr('refY', 0)
    .attr('markerWidth', 8)
    .attr('markerHeight', 8)
    .attr('orient', 'auto')
    .append('path')
    .attr('d', 'M0,-5L10,0L0,5')
    .attr('fill', '#666')
  
  // 创建连线
  const link = svg.append('g')
    .selectAll('line')
    .data(data.links)
    .enter().append('line')
    .attr('stroke', '#999')
    .attr('stroke-width', 2)
    .attr('marker-end', 'url(#end)')
  
  // 创建节点组
  const node = svg.append('g')
    .selectAll('g')
    .data(data.nodes)
    .enter().append('g')
    .call(d3.drag()
      .on('start', dragstarted)
      .on('drag', dragged)
      .on('end', dragended))
  
  // 添加节点圆圈
  node.append('circle')
    .attr('r', 20)
    .attr('fill', '#1890ff')
    .attr('stroke', '#fff')
    .attr('stroke-width', 2)
  
  // 添加节点标签
  node.append('text')
    .text(d => d.id)
    .attr('font-size', 12)
    .attr('text-anchor', 'middle')
    .attr('dy', 4)
    .attr('fill', '#fff')
  
  // 更新力导向图
  simulation.on('tick', () => {
    link
      .attr('x1', d => d.source.x)
      .attr('y1', d => d.source.y)
      .attr('x2', d => d.target.x)
      .attr('y2', d => d.target.y)
    
    node.attr('transform', d => `translate(${d.x},${d.y})`)
  })
  
  // 添加缩放功能
  const zoom = d3.zoom()
    .scaleExtent([0.1, 4])
    .on('zoom', (event) => {
      svg.selectAll('g').attr('transform', event.transform)
    })
  
  svg.call(zoom)
  
  function dragstarted(event) {
    if (!event.active) simulation.alphaTarget(0.3).restart()
    event.subject.fx = event.subject.x
    event.subject.fy = event.subject.y
  }
  
  function dragged(event) {
    event.subject.fx = event.x
    event.subject.fy = event.y
  }
  
  function dragended(event) {
    if (!event.active) simulation.alphaTarget(0)
    event.subject.fx = null
    event.subject.fy = null
  }
}

const parseStreamResponse = (text) => {
  try {
    // 尝试解析为 JSON
    return JSON.parse(text)
  } catch (e) {
    // 如果不是有效的 JSON，尝试提取可能的 JSON 部分
    const jsonMatch = text.match(/\{[\s\S]*\}/)
    if (jsonMatch) {
      try {
        return JSON.parse(jsonMatch[0])
      } catch (e) {
        console.error('Failed to parse JSON from stream:', e)
        return null
      }
    }
    return null
  }
}

const handleSubmit = async () => {
  if (!inputValue.value) {
    message.warning('请输入内容')
    return
  }
  
  loading.value = true
  answerText.value = ''
  streamLogs.value = []
  showStreamLogs.value = true
  graphData.value = null
  
  try {
    const response = await fetch('http://192.168.5.6/v1/chat-messages', {
      method: 'POST',
      headers: {
        'Authorization': 'Bearer app-aQmrylCiemu3Ksg2xrLexgG5',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        inputs: {},
        query: inputValue.value,
        response_mode: 'streaming',
        conversation_id: '',
        user: 'abc-123'
      })
    })

    const reader = response.body.getReader()
    const decoder = new TextDecoder()
    let buffer = ''

    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      const chunk = decoder.decode(value, { stream: true })
      buffer += chunk

      // 处理 SSE 格式的数据
      const lines = buffer.split('\n')
      buffer = lines.pop() || '' // 保留最后一个不完整的行

      for (const line of lines) {
        if (line.startsWith('data: ')) {
          const data = line.slice(6) // 移除 'data: ' 前缀
          if (data.trim() === '') continue // 跳过空行
          
          try {
            const jsonData = JSON.parse(data)
            addLog(jsonData.event || 'unknown', jsonData)
            
            if (jsonData.event === 'agent_message' && jsonData.answer) {
              answerText.value += jsonData.answer
            }
            
            // 检查是否包含图数据
            if (jsonData.nodes && jsonData.links) {
              graphData.value = jsonData
              drawGraph(jsonData)
              // 当图数据生成完成后，隐藏流式日志
              setTimeout(() => {
                showStreamLogs.value = false
              }, 1000)
            }
          } catch (e) {
            console.error('Failed to parse SSE data:', e)
          }
        }
      }
    }
  } catch (error) {
    message.error('请求失败：' + error.message)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="app-container">
    <header class="app-header">
      <h1>AI 知识图谱可视化</h1>
    </header>
    <main class="app-main">
      <div class="content-wrapper">
        <div class="left-panel">
          <div class="input-section">
            <Input.TextArea
              v-model:value="inputValue"
              placeholder="请输入您的问题..."
              :auto-size="{ minRows: 3, maxRows: 6 }"
              :disabled="loading"
              class="input-area"
            />
            <Button 
              type="primary" 
              @click="handleSubmit" 
              class="submit-btn"
              :loading="loading"
            >
              生成图谱
            </Button>
          </div>
          
          <div v-if="showStreamLogs" class="stream-section">
            <div class="stream-container">
              <div v-for="(log, index) in streamLogs" :key="index" class="stream-item">
                <div class="stream-header">
                  <span class="stream-type" :class="log.type">{{ log.type }}</span>
                  <span class="stream-timestamp">{{ log.timestamp }}</span>
                </div>
                <pre class="stream-content">{{ JSON.stringify(log.data, null, 2) }}</pre>
              </div>
            </div>
          </div>
        </div>

        <div class="right-panel">
          <div v-if="graphData" class="graph-section">
            <div class="graph-container">
              <svg ref="svgRef"></svg>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<style>
.app-container {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #f8f9fa;
}

.app-header {
  padding: 1rem 1.5rem;
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  border-bottom: 1px solid #e8e8e8;
  position: sticky;
  top: 0;
  z-index: 100;
}

.app-header h1 {
  margin: 0;
  color: #1a1a1a;
  font-size: 1.5rem;
  font-weight: 600;
  background: linear-gradient(45deg, #1890ff, #722ed1);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.app-main {
  flex: 1;
  padding: 1rem;
  overflow: hidden;
}

.content-wrapper {
  height: 100%;
  display: flex;
  gap: 1rem;
  max-width: 1600px;
  margin: 0 auto;
}

.left-panel {
  width: 400px;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  overflow-y: auto;
}

.right-panel {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

.input-section {
  background: #fff;
  padding: 1.5rem;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.input-area {
  width: 100%;
  margin-bottom: 1rem;
}

.input-area textarea {
  border-radius: 8px;
  font-size: 1rem;
  line-height: 1.6;
  padding: 1rem;
}

.submit-btn {
  width: 100%;
  height: 40px;
  font-size: 1rem;
  border-radius: 8px;
  background: linear-gradient(45deg, #1890ff, #722ed1);
  border: none;
}

.submit-btn:hover {
  opacity: 0.9;
  transform: translateY(-1px);
}

.stream-section {
  background: #fff;
  padding: 1.5rem;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.stream-container {
  flex: 1;
  overflow-y: auto;
  padding-right: 1rem;
}

.stream-container::-webkit-scrollbar {
  width: 6px;
}

.stream-container::-webkit-scrollbar-track {
  background: #f5f5f5;
  border-radius: 3px;
}

.stream-container::-webkit-scrollbar-thumb {
  background-color: #d9d9d9;
  border-radius: 3px;
}

.stream-item {
  margin-bottom: 1rem;
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.stream-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stream-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.stream-type {
  padding: 0.3rem 0.6rem;
  border-radius: 6px;
  font-size: 0.85rem;
  font-weight: 500;
}

.stream-type.agent_message {
  background: #e6f7ff;
  color: #1890ff;
}

.stream-type.agent_thought {
  background: #f6ffed;
  color: #52c41a;
}

.stream-timestamp {
  color: #8c8c8c;
  font-size: 0.8rem;
}

.stream-content {
  margin: 0;
  font-family: 'Fira Code', monospace;
  font-size: 0.85rem;
  line-height: 1.6;
  color: #2c3e50;
  white-space: pre-wrap;
  word-break: break-all;
}

.graph-section {
  background: #fff;
  padding: 1.5rem;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  height: 100%;
  display: flex;
  flex-direction: column;
}

.graph-container {
  flex: 1;
  background: #f8f9fa;
  border-radius: 8px;
  overflow: hidden;
}

svg {
  width: 100%;
  height: 100%;
}

@media (max-width: 768px) {
  .content-wrapper {
    flex-direction: column;
  }
  
  .left-panel {
    width: 100%;
    height: 50vh;
  }
  
  .right-panel {
    height: 50vh;
  }
}
</style>
