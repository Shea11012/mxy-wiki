---
created: 2025-06-27
tags: [js/fetch]
updated: 2025-06-27
---

## 请求配置

POST 请求示例
```js
fetch('https://api.example.com/data', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer your_token_here'
  },
  body: JSON.stringify({
    name: 'John Doe',
    email: 'john@example.com'
  })
})
.then(response => response.json())
.then(data => console.log('Success:', data))
.catch(error => console.error('Error:', error));
```

发送 FormData
```js
const formData = new FormData();
formData.append('username', 'john');
formData.append('avatar', fileInput.files[0]);

fetch('https://api.example.com/profile', {
  method: 'POST',
  body: formData
})
.then(response => response.json())
.then(result => console.log('Success:', result))
.catch(error => console.error('Error:', error));
```

## 流式响应

JSON 数据
```js
fetch('https://api.example.com/data.json')
  .then(response => response.json())
  .then(data => console.log(data));
```

文本数据
```js
fetch('https://api.example.com/data.txt')
  .then(response => response.text())
  .then(text => console.log(text));
```

二进制数据
```js
fetch('https://api.example.com/image.png')
  .then(response => response.blob())
  .then(blob => {
    const imgUrl = URL.createObjectURL(blob);
    document.getElementById('myImg').src = imgUrl;
  });
```

ArrayBuffer 原始二进制数据
```js
fetch('https://api.example.com/data.bin')
  .then(response => response.arrayBuffer())
  .then(buffer => {
    // 处理二进制数据
    const view = new DataView(buffer);
    console.log(view.getInt32(0));
  });
```

FormData 表单数据
```js
fetch('https://api.example.com/form')
  .then(response => response.formData())
  .then(formData => {
    console.log(formData.get('username'));
  });
```

分块读取文本数据
```js
fetch('https://api.example.com/large-text-file.txt')
  .then(response => {
    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    let result = '';
    
    function readChunk({ done, value }) {
      if (done) {
        console.log('完整内容:', result);
        return;
      }
      
      result += decoder.decode(value, { stream: true });
      return reader.read().then(readChunk);
    }
    
    return reader.read().then(readChunk);
  });
```

处理大型 JSON 数据流
```js
async function processLargeJsonStream(url) {
  const response = await fetch(url);
  const reader = response.body.getReader();
  const decoder = new TextDecoder();
  let buffer = '';
  
  while (true) {
    const { done, value } = await reader.read();
    if (done) break;
    
    buffer += decoder.decode(value, { stream: true });
    
    // 处理缓冲区中完整的JSON对象
    let boundary;
    while ((boundary = buffer.indexOf('}\n{')) !== -1) {
      const completeJson = buffer.substring(0, boundary + 1);
      buffer = buffer.substring(boundary + 2);
      
      try {
        const data = JSON.parse(completeJson);
        console.log('处理数据:', data);
      } catch (e) {
        console.error('JSON解析错误:', e);
      }
    }
  }
  
  // 处理剩余数据
  if (buffer) {
    try {
      const data = JSON.parse(buffer);
      console.log('最后的数据:', data);
    } catch (e) {
      console.error('最后JSON解析错误:', e);
    }
  }
}

processLargeJsonStream('https://api.example.com/streaming-data');
```

## 中断请求

```js
const controller = new AbortController();
const signal = controller.signal;

// 设置5秒超时
const timeoutId = setTimeout(() => controller.abort(), 5000);

fetch('https://api.example.com/data', { signal })
  .then(response => response.json())
  .then(data => {
    clearTimeout(timeoutId);
    console.log(data);
  })
  .catch(error => {
    if (error.name === 'AbortError') {
      console.log('请求被中止');
    } else {
      console.error('请求出错:', error);
    }
  });
```

## 并行请求

```js
Promise.all([
  fetch('https://api.example.com/users'),
  fetch('https://api.example.com/posts')
])
.then(async ([usersResponse, postsResponse]) => {
  const users = await usersResponse.json();
  const posts = await postsResponse.json();
  return { users, posts };
})
.then(data => console.log('Combined data:', data))
.catch(error => console.error('Error in one of the requests:', error));
```

## 错误处理

```js
async function fetchWithErrorHandling(url, options = {}) {
  try {
    const response = await fetch(url, options);
    
    if (!response.ok) {
      // 尝试解析错误信息
      let errorInfo;
      try {
        errorInfo = await response.json();
      } catch (e) {
        errorInfo = await response.text();
      }
      
      // 根据状态码抛出不同类型的错误
      if (response.status >= 400 && response.status < 500) {
        throw new ClientError(response.status, errorInfo);
      } else if (response.status >= 500) {
        throw new ServerError(response.status, errorInfo);
      }
    }
    
    return await response.json();
  } catch (error) {
    if (error instanceof DOMException && error.name === 'AbortError') {
      console.log('请求被用户中止');
      throw error;
    }
    
    console.error('请求失败:', error);
    throw error;
  }
}

class ClientError extends Error {
  constructor(status, message) {
    super(`客户端错误: ${status} - ${message}`);
    this.status = status;
  }
}

class ServerError extends Error {
  constructor(status, message) {
    super(`服务器错误: ${status} - ${message}`);
    this.status = status;
  }
}

// 使用示例
fetchWithErrorHandling('https://api.example.com/data')
  .then(data => console.log(data))
  .catch(error => console.error('处理后的错误:', error));
```

