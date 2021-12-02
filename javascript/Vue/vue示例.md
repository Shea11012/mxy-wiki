### 在 vue 内创建一个弹窗

```js
import Vue from 'vue';
export default function create(component,props) {
	const vm = new Vue({ // 先创建vue一个实例
		render(h) {
			return h(component,{props})  // 使用 render 提供的 h 方法渲染组件
		}
	}).$mount(); 	// 更新

	// 在 vm 实例上挂载了一个组件实例，可以通过 $children 获取该实例，因为只挂载了一个可以通过指定索引获取
	const comp = vm.$children[0];
	// 将该实例的 dom 挂载到 body 上
	document.body.appendChild(vm.$el)
	// 给组件实例添加一个清理函数，在使用完后销毁该实例，避免内存溢出
	comp.remove = () => {
		document.body.removeChild(vm.$el);
		vm.$destroy();
	}

	return comp;
}

// 使用的时候保证组件上定义了需要的 props 
```

