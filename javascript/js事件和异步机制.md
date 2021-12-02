# JavaScript事件和异步机制

[js事件和异步机制](https://www.digitalocean.com/community/tutorials/understanding-the-event-loop-callbacks-promises-and-async-await-in-javascript)

```javascript
function first() {
  console.log(1);
}

function second() {
  console.log(2);
}

function third() {
  console.log(3);
}

first();
second();
third();

// 输出结果
1
2
3
```

上面的代码是按照执行顺序输出的，但是当使用了异步的web API 的时候，这个输出结果就会不一样了

```javascript
function first() {
  console.log(1);
}

function second() {
  setTimeout(() => {
	  console.log(2);
  })
}

function third() {
  console.log(3);
}

first();
second();
third();

// 输出结果
1
3
2
```

导致这种结果的原因是因为JavaScript内的 **stack**。

### Stack

stack，内存放了当前正在运行和待运行的函数，浏览器处理规则如下：

代码片段1

- `first` 入栈，运行 `first` 输出 `1`，将 `first` 移出栈
- `second` 入栈，运行 `second` 输出 `2` ，将 `second` 移出栈
- `third` 入栈，运行 `third` 输出 `3` ，将 `third` 移出栈

代码片段2

- `first` 入栈，运行 `first`  输出 `1`，将 `first` 移出栈
- `second` 入栈，运行 `second` 
  - 将 `setTimeout` 入栈，运行 `setTimeout` 开启一个计时器且将匿名函数加入到 `queue`中，将 `setTimeout` 移出栈
- `second` 移出栈
- `third` 入栈，运行 `third` 输出 `3`，`third` 移出栈
- 最后event loop 检查 `queue` 是否有待执行函数，将来自 `setTimeout` 的匿名函数入栈运行，输出 `2`，再移出栈。

### Queue

queue，可以被看做是消息队列或者任务队列，是函数等待区域无论栈是否空置。event loop 将从队头开始检查并执行。如 `setTimeout` 函数，它会立即将匿名函数放置进队列中，**就算 setTimeout 设置的时间是 0 秒，也并不意味着它会在 0 秒后执行，而是会在该时间将匿名函数添加到stack中。**

> Note：还有另外一种队列叫做 Job queue 工作队列或者 microtask queue 微任务队列，这种队列是为 promises 准备的。Microtask 如 promises 比 setTimeout 这种 macrotask 具有更高的优先级。

在为了确保使用setTimeout后还能让程序按照正常顺序输出，可以使用回调的方式保证输出的顺序

```javascript
function first() {
  console.log(1);
}

function second(callback) {
  setTimeout(() => {
    console.log(2);
    callback();
  });
}

function third() {
  console.log(3);
}

first();
second(third);

// 输出结果
1
2
3
```

上面代码通过回调的方式，很容易造成回调地狱，使得代码难以读懂。

### Promise

promise 异步函数，具有三种状态：

- Pending：在 resolved 或 rejected 之前初始化状态
- Fulfilled：成功的操作，promise 已经被 resolved
- Rejected：失败操作，promise 已经被 rejected

### Async await

async 函数使得异步代码看上去就像同步代码一样，async 内必须配合 await 使用

```javascript
async function getUser() {
  const data = await http.get(url);
  return data;
}
```

