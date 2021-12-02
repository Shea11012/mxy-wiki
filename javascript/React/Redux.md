[redux 文档](https://redux.js.org/)

创建一个 store，store 管理全局的 state

```react
const reducer = (state,action) => {
    // state 作为状态管理，不可直接操作
    // action 是一个对象，必须带有 type 属性，通过 type 判断要操作的方法
    // reducer 必须返回 state
    switch(action.type){
        case "LOGIN":
            return Object.assign({},state,{}); // 这是一种返回 state 的方法
        case "LOGOUT":
            return {...state,{}};	// 必须支持 es6
    	default：
        	return state；
    }
}
const store = Redux.createStore(reducer);

store.getState(); 	// 得到当前的 state
store.dispatch(Object);	// 接受一个对象
store.listener(function); // 在每次 dispatch 时，都会触发这个 listener

// 组合多个 reducers 变成一个 reducer
const reducer1 = (state,action)=>{};
const reducer2 = (state,action)=>{};
const rootReducer = Redux.combineReducers({
    one:reducer1,
    two:reducer2,
});
Redux.createStore(rootReducer);
```

#### 使用 middleware 处理 async actions

```react
const store = Redux.createStore(
  asyncDataReducer, // 异步处理 reducer
  Redux.applyMiddleware(ReduxThunk.thunk)
);
```

#### redux 和 react 相连

```react
// 使用 Provider 
import {Prodiver} from 'react-redux';
import {createStore} from 'react';
import {render} from 'react-dom';

const store = create(reducer);

render (
	<Provider store={store}>	// 通过 provider 包裹根组件，再传递一个 store,就可以使用 redux
        <App />
    </Provider>,
    document.getElementById('root')
);

// connect 方法
```

