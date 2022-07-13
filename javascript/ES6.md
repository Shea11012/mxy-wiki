---
date created: 2021-12-03 20:20
date modified: 2022-04-29 20:15
title: ES6
---
# ES6

### 变量的解构赋值

```javascript
let [x = 1] = [undefined];	//x = 1
let [x = 1] = [null];		//x = null
```

如果是一个表达式，那么这个表达式是惰性求值：只有在会用到的时候，才会求值。

```javascript
function f(){
    console.log('aaa');
}
let [x = f()] = [1];	//这里因为x能取到值，函数f不会执行

//等价于下面
let x;
if([1][0] === undefined){
    x = f();
}else{
    x = [1][0];
}
```

对象解构

```javascript
let obj = {p:['hello',{y:'world'}]};
let {p:[x,{y}]} = obj;
x // hello
y //world
```

### 字符串

JavaScript 内部，字符以 UTF-16 格式存储，每个字符固定为 2 个字节，大于 `0xFFFF` 的字符被认为是两个字符



### 正则扩展

先行断言：x 只有在 y 前面才匹配，必须写成 `/x(?=y)/` ，例： `/\d+(?=%)/.exec('100% of US presidents have been male')` ，结果： `['100']`

先行否定断言：x 只有不在 y 前面才匹配，必须写成 `/x(?!y)/`，例： `/\d+(?!%)/.exec('that's all 44 of them')`， 结果： `['44']`

后行断言：x 只有在 y 后面才匹配，必须写成 `/(?<=y)x/` ，例： `/(?<=\$)\d+/.exec('Benjamin Franklin is on the $100 bill')` ，结果： `['100']`

后行否定断言：x 只有不在 y 后面才匹配，必须写成 `/(?<!y)x/` ，例： `/(?<!\$)\d+/.exec('it’s is worth about €90')` ，结果： `['90']`

注意：断言括号之中部分，是不计入返回结果的。

后行断言的实现，是先匹配 `/(?<=y)x/` 的 `x` ，然后在回到左边，匹配 `y` 的部分

```javascript
/(?<=(\d+)(\d+))$/.exec('1053') // ["", "1", "053"]
/^(\d+)(\d+)$/.exec('1053') // ["1053", "105", "3"]
```

后行断言，由于执行顺序是从右到左，第二个括号是贪婪模式，第一个括号只能捕获一个字符

非后行断言，第一个括号是贪婪模式，第二个括号只能捕获一个字符

后行断言的反斜杠引用，与通常顺序相反，必须放在对应的那个括号之前

```javascript
/(?<=(o)d\1)r/.exec('hodor')  // null
/(?<=\1d(o))r/.exec('hodor')  // ["r", "o"]
```



**Unicode 属性类**

`\p{.....}` 表示匹配，`\P{......}` 表示反向匹配

```javascript
// 匹配所有空格
\p{White_Space}

// 匹配各种文字的所有字母，等同于 Unicode 版的 \w
[\p{Alphabetic}\p{Mark}\p{Decimal_Number}\p{Connector_Punctuation}\p{Join_Control}]

// 匹配各种文字的所有非字母的字符，等同于 Unicode 版的 \W
[^\p{Alphabetic}\p{Mark}\p{Decimal_Number}\p{Connector_Punctuation}\p{Join_Control}]

// 匹配 Emoji
/\p{Emoji_Modifier_Base}\p{Emoji_Modifier}?|\p{Emoji_Presentation}|\p{Emoji}\uFE0F/gu

// 匹配所有的箭头字符
const regexArrows = /^\p{Block=Arrows}+$/u;
regexArrows.test('←↑→↓↔↕↖↗↘↙⇏⇐⇑⇒⇓⇔⇕⇖⇗⇘⇙⇧⇩') // true
```

**具名组匹配**：允许给每一个组匹配定一个名字，便于引用，命名方式 `?<name>`

```javascript
const RE_DATE = /(?<year>\d{4})-(?<month>\d{2})-(?<day>\d{2})/;

const matchObj = RE_DATE.exec('1999-12-31');
const year = matchObj.groups.year; // 1999
const month = matchObj.groups.month; // 12
const day = matchObj.groups.day; // 31
```



### 函数扩展

#### 函数结构赋值使用

```javascript
function foo({x,y=5}) {
    console.log(x,y);
}
foo({}) // undefined 5
foo({x:1}) // 1 5
foo({x:1,y:2}) // 1 2
foo() // TypeError:Cannot read property 'x' of undefined 
// 会报 Type Error错误是因为只使用了对象的解构赋值默认值，没有使用函数参数的默认值。只用当函数 foo 的参数是一个对象时，变量 x 和 y 才会通过结构赋值生成。

function foo({x,y=5} = {}) {
    console.log(x,y);
}
foo() // undefined 5 

```



`reset` 参数写法 `function add(...values)` 

箭头函数注意点：

1. 函数体内的 this 对象，就是定义时所在的对象，而不是使用时所在的对象
2. 不可以当作构造函数，也就是说，不可以使用 new 命令，否则会抛出一个错误
3. 不可以使用 arguments 对象，该对象在函数体内不存在。如果要用，可以用 reset 参数代替
4. 不可以使用 yield 命令，因此箭头函数不能用作 Generator 函数。

**双冒号运算符：**双冒号左边是一个对象，右边是一个函数。该运算符会自动将左边的对象，作为上下文环境（即 this 对象），绑定到右边的函数上面。

```javascript
foo::bar;
//等同于
bar.bind(foo);

foo::bar(...arguments);
//等同于
bar.apply(foo,arguments);
```

**尾调用、尾调用优化、尾递归**

尾调用：指某个函数的最后一步是调用另一个函数，不一定出现在函数尾部，只要是最后一部操作即可。

```javascript
function f(x)
{
    if(x>0){
        return m(x);
    }
    return n(x);
}
```

尾调用优化：只保留内层函数的调用帧，如果所有函数都是尾调用，那么完全可以做到每次执行，调用帧只有一项，这将大大节省内存。

注意：只有不再用到外层函数的内部变量，内层函数的调用帧才会取代外层函数的调用帧，否则就无法进行尾调用优化

```javascript
function addOne(a)
{
    var one = 1;
    function inner(b)
    {
        return b+one;
    }
    return inner(a);
}
// 因为内层函数 inner 用到外层函数 addone 的内部变量 one ， 上面的函数不会进行尾调用优化
```

尾递归：递归非常耗费内存，因为需要同时保存成千上百个调用帧，很容易发生 栈溢出 错误（stack overflow）。对于尾递归来说，由于只存在一个调用帧，所以永远不会发生 栈溢出 错误。

```javascript
// 未使用尾递归
function factorial(n)
{
    if (n === 1) return 1;
    return n * factorial(n - 1);
}

// 使用尾递归
function factorial(n,total)
{
    if (n === 1) return total;
    return factorial(n-1,n*total);
}
```

尾递归的实现，往往需要该写递归函数，确保最后一步只调用自身。需要把所有用到的内部变量该写成函数的参数

函数式编程有一个概念，叫做柯里化（currying），意思是将多参数函数转换成单参数的形式。

```javascript
function currying(fn,n)
{
    return function(m){
        return fn.call(this,m,n);
    };
}

function tailFactorial(n,total)
{
    if(n === 1) return total;
    return tailFactorial(n-1,n*total);
}

const facorial = currying(tailFactorial,1);
factorial(5); //120
```

尾调用仅在严格模式有效



### 数组的扩展

**扩展运算符：**`...`

复制数组： `const a1 = [1,2]; const a2 = [...a1];` ，这是浅拷贝

合并数组： `[...arr1,...arr2,...arr3]`



### 对象扩展

#### in 

in 关键字判断一个对象是否包含某个属性

```js
const a = {
  q: 3
  w: 4
}

console.log("q" in a) // true

// 同时 in 也可以使用在数组上
const arr = [1,2,3]
console.log(3 in arr) // false 判断的索引 3 位置上是否有值
```



对象属性名允许使用表达式

##### 属性的可枚举性和遍历

对象的每个属性都有一个描述对象（Descriptor），用来控制该属性的行为，`Object.getOwnPropertyDescriptor` 方法获取该属性的描述对象

```javascript
let obj = { foo: 123 };
Object.getOwnPropertyDescriptor(obj,'foo');
// {
//	value:123,
//	writable:true,
//	enumerable:true,	// 可枚举性
//	configurable:true,
// }
```

如果 `enumerable：false` 会有四个操作忽略

-  `for...in` ：只遍历对象自身的和继承的可枚举的属性
- `Object.keys()` ：返回对象自身的所有可枚举属性和键名
- `JSON.stringify()` ：只串化对象自身的可枚举属性
- `Object.assign()` ：只拷贝对象自身的可枚举属性

##### 属性的遍历

1. `for...in`：遍历自身和继承的可枚举属性（不含 Symbol 属性）
2. `Object.keys(obj)` ：返回一个数组，包括对象自身的（不含继承的）所有可枚举属性（不含 Symbol 属性）的键名
3. `Object.getOwnPropertyNames(obj)`：返回一个数组，包含对象自身的所有属性（不含 Symbol 属性，包括不可枚举属性）的键名
4. `Object.getOwnPropertySymbols(obj)`：返回一个数组，包含对象自身的所有 Symbol 属性的键名
5. `Reflect.ownKeys(obj)` ：返回一个数组，包含对象自身的所有键名

**以上遍历对象的键名，遵守同样的属性遍历规则：**

- 先遍历所有数值键，按照数值升序排列
- 再遍历所有字符串键，按照加入时间升序排列
- 最后遍历所有 Symbol 键，按照加入时间升序排列

##### 解构赋值

对象的解构赋值相当于将目标对象自身的所有可遍历的 enumerable ，分配到指定的对象上面。

>解构赋值要求等号右边是一个对象
>
>解构赋值是浅拷贝
>
>扩展运算符的解构赋值，不能复制继承自原型对象的属性

### Symbol

唯一的符号，用于取代会发生冲突的 key 和值定义

```js
const s = Symbol()

const areaType = {
  square: Symbol(),
  cicrle: Symbol(),
}

function getArea(type) {
  switch(type) {
    case areaType.square:
      return 1
    case areaType.cicrle:
      return 2
  }
}
```



### Set 和 Map 数据结构

#### Set

Set 实例的属性和方法

实例：

- `Set.prototype.constructor `：构造函数，默认就是 Set 函数
- `Set.prototype.size` ：返回 Set 实例的成员总数

方法：分为操作方法和遍历方法

- `add(value)` ：添加某个值，返回 Set 结构本身
- `delete(value)` ：删除某个值，返回一个布尔值，表示删除是否成功
- `has(value)` ：返回一个布尔值，表示该值是否为 Set 的成员
- `clear()` ：清除所有成员，没有返回值
- `keys()`  ：返回键名的遍历器
- `values()` ： 返回键值的遍历器
- `entries()` ：返回键值对的遍历器
- `forEach()` ： 使用回调函数遍历每个成员

Set 结构没有键名，只有键值（或者说键名和键值是同一个值），所以 keys 和 values 方法一致

#### WeakSet

WeakSet 成员只能是对象，WeakSet 中的对象都是弱引用，即垃圾回收机制不考虑 WeakSet 对该对象的引用，如果其他对象都不再引用该对象，那么垃圾回收机制会自动回收该对象所占用的内存，不考虑该对象还存在于 WeakSet 之中。

因为垃圾回收机制依赖引用计数，如果一个值的引用次数不为 0，垃圾回收机制就不会释放这块内存。结束使用该值之后，有时会忘记取消引用，导致没存无法释放，进而可能会引发内存泄漏。WeakSet 里面的引用，都不计入垃圾回收机制，所以就不存在这个问题。因此，WeakSet 适合临时存放一组对象，以及存放跟对象绑定的信息。只要这些对象在外部消失，它在 WeakSet 里面的引用就会自动消失。 WeakSet 内部有多少个成员，取决于垃圾回收机制有没有运行，运行前后很可能成员个数不一样，因此 WeakSet 不可遍历。

#### Map

Map 数据结构类似对象，是键值对的集合，但是键不限于字符串，各种类型的值包括对象都可以当做键

#### WeakMap

### Proxy

用于修改目标对象的默认操作

`new Proxy(target,handler)`

target : 表示所要拦截的目标对象

handler ：定制拦截行为



### Reflect

- 将 Object 内的一些方法移到了 reflect 上
- 修改某些 Object 返回的结果

```js
// es5
try {
  Object.defineProperty()
} catch (e) {
  
}

// es6
if (Reflect.defineProperty()) {
  
}
```



### Promise 对象

Promise 是异步编程的一种解决方案，它有三种状态，pending（进行中），resolved（已完成），rejected（已失败）。

当 Promise 的状态发生变化时会执行相应的方法，并且状态一旦改变，就无法再次改变状态。

```javascript
//实例化的promise对象会立即执行
let promise = new Promise((resolve,reject)=>{
    if(success){
        resolve(a)	//pending ---> resolved
    }else{
        reject(err)	//pending ---> rejectd
    }
})

function promise(){
    return new Promise(function(resolve,reject){
        if(success){
            resolve(a)
        }else{
            reject(err)
        }
    })
}
```

`.then()` 方法是 promise 原型链上的方法，它包含两个参数方法，分别是 resolve 和 reject

```javascript
promise.then(
	()=>{console.log('this is success callback')},
    ()=>{console.log('this is fail callback')}
)
```

`.catch()` 是捕获 Promise 错误，Promise 的抛错具有冒泡性质，能够不断传递，所以建议不要使用 then() 的 reject 回调，而是统一使用 catch() 来处理错误

```javascript
promise.then(
	()=>{console.log('this is success callback')}
).catch(
	(err)=>{console.log(err)}
)
//catch()也可以抛出错误，抛出的错误会在下一个catch中捕获处理，因此可以再添加catch()
```

#### Promise.resolve() 和 Promise.reject()

Promise.resolve()：

- 参数是 Promise，原样返回
- 参数带有 then 方法，转换为 Promise 后立即执行 then 方法
- 参数不带 then 方法，不是对象或没有参数，返回 resolve 状态的 Promise

Promise.reject()：会直接返回 rejected 状态的 promise

#### Promise.all()

参数为 Promise 对象数组，如果不是 promise 对象，会通过 Promise.resolve() 方法转换

```javascript
var promise = Promise.all([p1,p2,p3])
promise.then(
......
).catch(
.....
)
//当p1,p2,p3的状态都变成resolved，promise才会变成resolved，并调用then()，但只要有一个变成rejected状态，promise就会立刻变成rejected状态
```

#### Promise.race()

```javascript
var promise = Promise.race([p1,p2,p3])
promise.then(
......
).catch(
......
)
//参数与Promise.all()相同，不同的是参数中的p1,p2,p3只要有一个改变状态，promise就会立刻变成相同的状态并执行对应的回调
```



### Iterator 和 for ..... of 循环

ES6 规定，默认的 Iterator 接口部署在数据结构的 `Symbol.iterator` 属性。一个数据结构只要具有 `Symbol.iterator` 属性，就认为是可遍历的

```javascript
class RangeIterator
{
    constructor(star,stop)
    {
        this.value = start;
        this.stop = stop;
    }
    [Symbol.iterator](){return this;}
    
    next(){
        let value = this.value;
        if(value < this.stop){
            this.value++;
            return {done:false,value:value};
        }
        return {done:true,value:undefined};
    }
}

function range(start,stop)
{
    return new RangeIterator(start,stop);
}

for(let value of range(0,3)){
    console.log(value);	// 0,1,2
}
```

不具有 Iterator 接口的对象或数组，可以使用 `Array.from` 方法将其转为数组。

```javascript
let arrayLike = {length:2,0:'a',1:'b'};
for(let x of Array.from(arrayLike)){
    console.log(x);
}
```

`for .... in` 循环缺点：

- 数组的键名是数字，但是 for....in 循环是以字符串作为键名 0，1，2 等。
- for...in 循环不仅遍历数字键名，还会遍历手动添加的其他键，甚至包括原型链上的键。
- 某些情况下，for....in 循环会以任意顺序遍历键名。
- 主要是为了遍历对象而设计的，不使用遍历数组

`for...of`

- 有着和 for...in 一样的简洁写法，但是没有 for...in 的那些缺点。
- 不同于 foreach 方法，可以与 break、continue 和 return 配合使用
- 提供了遍历所有数据结构的统一操作接口



### Generator 语法和异步应用

generator 实现了 es6 的协程，最大的特点就是可以交出函数的执行权（即暂停执行）。异步操作需要暂停的地方，都用 `yield` 语句实现。

Generator 函数是一个普通函数，有两个特征。

一、`function` 关键字与函数名之间有一个星号 `*`

二、函数体内部使用 `yield` 表达式，定义不同的内部状态。

##### next 方法参数

```javascript
function* foo(x)
{
    let y = 2 * (yield (x + 1));
    let z = yield (y / 3);
    return (x + y + z);
}

let a = foo(5);
a.next();	// Object(value:6,done:false)
a.next();	// Object(value:NaN,done:false)
a.next();	// Object(value:NaN,done:true)

let b = foo(5);
b.next();	// {value:6,done:false}
b.next(12)	// {value:8,done:false}
b.next(13)	// {value:42,done:true}
```

当 next 不带参数时，导致 y 的值等于 `2 * undefined` 即 NaN，除以 3 以后还是 NaN，因此返回对象的 value 属性也等于 NaN。第三次运行 next 方法的时候不带参数，所以 Z 等于 undefined ，返回对象的 value 属性等于 `5 + NaN + undefined` 即 NaN。

如果 Next 方法提供参数，返回结果就完全不一样了，第一次调用 Next 返回 `x + 1` 的值 6，第二次调用 next 方法，将上一次 yield 表达式的值设为 12，因此 y 等于 24，返回 `y / 3` 的值，第三次调用 next 方法，将上一次 yield 表达式的值设为 13，因此 z 等于 13，此时，x 等于 5，y 等于 24，z 等于 13，所以 return 语句的值等于 42。

第一次使用 next 方法是，传递参数无效，V8 引擎直接忽略第一次使用 next 方法时的参数



### Async 函数

async 函数是 Generator 函数语法糖

```JavaScript
const as = async function(){
    const f1 = await func1();
    const f2 = await func2();
};
```

async 返回 promise 对象，可以用 then 指定下一步操作。

`for async (const x of items)` 遍历异步的接口，也可以遍历同步接口



### Class 语法

类的构造方法 `constructor`

类的方法之间不需要逗号分隔

在一个方法前，加上 `static` 关键字，表示该方法不会被实例继承，静态方法内的 `this` 指向类，不是指向实例

类的所有方法都实际定义在类的 `prototype` 属性

- getter 取值函数
- setter 存值函数

```javascript
class MyClass {
    constructor(){}
    
    get prop(){
        return 'getter';
    }
    
    set prop(val){
        this.prop = val;
    }
}
```

### class 继承

#### super

`super` 关键字表示父类的构造函数，用来新建父类的 `this` 对象

`super` 指向父类的原型对象，所以定义在父类实例上的方法或属性，无法通过 `super` 调用

`super` 作为对象，用在静态方法中，这时 `super` 将指向父类，而不是父类的原型对象

```javascript
class Parent {
    static myMethod(msg) {
        console.log('static',msg);
    }
    
    myMethod(msg) {
        console.log('instance',msg);
    }
}
class Child extends Parent {
    static myMethod(msg) {
        super.myMethod(msg);
    }
    
    myMethod(msg) {
        super.myMethod(msg);
    }
}
Child.myMethod(1); // static 1
let child = new Child();
child.myMethod(1);	// instance 1
```



子类必须在 `constructor` 方法中调用 `super` 方法，否则新建实例会报错，因为子类的 `this` 对象，必须通过父类的构造函数完成塑造，得到与父类同样的实例属性和方法

`Object.getPrototypeOf()` 判断一个类是否继承另一个类

##### 类的 `prototype` 属性和 `__proto__` 属性 

1. 子类的 `__proto__` 属性，表示构造函数的继承，总是指向父类
2. 子类 `prototype` 属性的 `__proto` 属性，表示方法的继承，总是指向父类的 `prototype` 属性

```javascript
class A {}
class B extend A {}
B.__proto__ === A // true
B.prototype.__proto__ === A.prototype // true
```

> 作为对象，子类的原型 `__proto__` 是父类
>
> 作为构造函数，子类的原型对象 `prototype` 是父类的原型对象

##### 实例的 `__proto__` 属性

子类的原型的原型，是父类的原型

```javascript
let p1 = new Point(2,3);
let p2 = new ColorPoint(2,3,'red');
p2.__proto__.__proto__ === p1.__proto__ //true
```

因此，通过子类的实例 `__proto__.__proto__` 属性，可以修改父类实例的方法

### Decorator 修饰器

```javascript
@testable
class MyTest {}
function testable(target) {
    target.isTestable = true;
}
MyTest.isTestable  // true
```

`@testable` 是一个修饰器，函数参数 `target` 是 `MyTest` 类本身

修饰器是对类的行为改变，是代码编译时发生的，而不是在运行时，修饰器的本质就是编译时执行的函数

#### 方法的修饰器

修饰器函数 `readonly` 可以接受三个参数

```javascript
class Person{
    @readonly
    name(){}
}

function readonly(target,name,descriptor){
    
}
// 第一个参数是类的原型对象
// 第二参数是要修饰的属性名
// 第三参数是该属性的描述对象
```

#### Trait 也是修饰器



### Module 语法

模块功能主要有两个命令构成：

- `export`：规定模块的对外接口

- `export default xx`：用于指定模块的默认输出，一个模块只能有一个默认输出
  ```javascript
  export const someVar = 123;
  export { someVar,someFunc };
  export { someVar as differentVarName }
  
  // 默认导出注意事项：
  // 在一个变量之前不需要使用 let/const/var;
  export default (someVar = 123);
  ```

- `import`：输入其他模块提供功能，在静态解析阶段执行，是一个模块中最早执行的

- `import()` ：运行时执行，不仅可以是模块，也可以是非模块的脚本，与 node 的 `require()` ，却别是前者是 异步加载，后者是同步加载

导出规则：当导入路径不是相对路径时，模块解析会模仿 node 模块解析策略。如：

`import * as foo from 'foo`，查找顺序如下：

1. `./node_modules/foo`
2. `../node_modules/foo`
3. `../../node_modules/foo`
4. 直到系统的根目录

### Module 加载实现

