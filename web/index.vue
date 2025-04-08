<template>
    <div>
      <button @click="sendRequest">发送请求</button>
      <pre>{{ responseText }}</pre>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        responseText: '' // 存储拼接后的完整文本
      };
    },
    methods: {
      async sendRequest() {
        const url = 'https://api.dify.ai/v1/chat-messages'; // 替换为您的 API 地址
        const headers = {
          'Authorization': 'Bearer app-xxxxxxxxxxxxxxxxxxxx', // 替换为您的 API Key
          'Content-Type': 'application/json'
        };
        const body = JSON.stringify({
          inputs: {},
          query: '对 192.168.5.2 进行安全验证',
          response_mode: 'streaming',
          conversation_id: '',
          user: 'user-123'
        });
  
        try {
          const response = await fetch(url, {
            method: 'POST',
            headers: headers,
            body: body
          });
  
          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
          }
  
          const reader = response.body.getReader();
          let fullAnswer = '';
  
          // 逐块读取流数据
          while (true) {
            const { done, value } = await reader.read();
            if (done) break; // 流结束
  
            // 将 Uint8Array 转换为字符串
            const chunk = new TextDecoder('utf-8').decode(value);
            const lines = chunk.split('\n');
  
            // 处理每一行
            for (const line of lines) {
              if (line.startsWith('data: ')) {
                const data = line.slice(6); // 去掉 "data: " 前缀
                try {
                  const parsed = JSON.parse(data); // 解析每块 JSON
                  fullAnswer += parsed.answer; // 拼接 answer 字段
                } catch (e) {
                  console.error('解析 JSON 失败:', e, '原始数据:', data);
                }
              }
            }
          }
  
          // 更新完整文本
          this.responseText = JSON.parse(fullAnswer); // 转换为对象
        } catch (error) {
          console.error('请求失败:', error);
          this.responseText = '请求失败，请检查控制台';
        }
      }
    }
  };
  </script>