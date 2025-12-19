import { Todo } from './types.js'; // 导入我们的类型

// 1. 导出我们需要操作的 DOM 元素
export const todoListElement = document.getElementById('todo-list');
export const newTodoForm = document.getElementById('new-todo-form') as HTMLFormElement | null;
export const todoNameInput = document.getElementById('todo-name-input') as HTMLInputElement | null;
export const todoDescriptionInput = document.getElementById('todo-description-input') as HTMLInputElement | null;

// 2. 导出主渲染函数
export function renderTodos(todos: Todo[]) {
    if (!todoListElement) return;
    todoListElement.innerHTML = ''; // 清空列表

    // 如果没有 todo，显示提示
    if (!todos || todos.length === 0) {
        const li = document.createElement('li');
        li.className = 'text-gray-500';
        li.textContent = '暂无待办事项';
        todoListElement.appendChild(li);
        return;
    }

    // 遍历 todos 并创建 HTML 元素
    todos.forEach(todo => {
        const li = document.createElement('li');
        li.className = 'p-4 bg-white rounded-lg shadow-md flex justify-between items-center';

        // 附加数据
        li.dataset.id = todo.id;
        li.dataset.name = todo.name;
        li.dataset.description = todo.description || "";
        li.dataset.completed = String(todo.completed);

        // --- 左侧 (复选框 + 文本) ---
        const leftDiv = document.createElement('div');
        leftDiv.className = 'flex items-center space-x-3';

        // 复选框
        const checkbox = document.createElement('input');
        checkbox.type = 'checkbox';
        checkbox.checked = todo.completed;
        checkbox.className = 'h-5 w-5 text-blue-600 rounded border-gray-300 focus:ring-blue-500';
        checkbox.dataset.action = 'toggle';

        // --- (新) 文本内容（垂直堆叠）---
        const textDiv = document.createElement('div');
        textDiv.className = 'flex flex-col';

        // 名称
        const nameSpan = document.createElement('span');
        nameSpan.textContent = todo.name;
        nameSpan.className = 'text-lg font-semibold';
        if (todo.completed) {
            nameSpan.classList.add('line-through', 'text-gray-400');
        }
        textDiv.appendChild(nameSpan);

        // --- (新) 显示 description ---
        // 仅当 description 存在时才显示
        if (todo.description) {
            const descriptionP = document.createElement('p');
            descriptionP.textContent = todo.description;
            descriptionP.className = 'text-sm text-gray-500'; // 灰色小号字体
            textDiv.appendChild(descriptionP);
        }

        // --- 组装左侧 ---
        leftDiv.appendChild(checkbox);
        leftDiv.appendChild(textDiv); // <-- 将 textDiv (而不是 nameSpan) 添加到 leftDiv

        // --- 右侧 (删除按钮) ---
        const deleteButton = document.createElement('button');
        deleteButton.textContent = '删除';
        deleteButton.className = 'bg-red-500 text-white px-3 py-1 rounded-md hover:bg-red-600 text-sm transition-colors';
        deleteButton.dataset.action = 'delete';

        // --- 组装 ---
        li.appendChild(leftDiv);
        li.appendChild(deleteButton);
        todoListElement.appendChild(li);
    });
}

// 3. 导出一个渲染错误的函数
export function renderError(message: string) {
    if (todoListElement) {
        todoListElement.innerHTML = `<li class="text-red-500">${message}</li>`;
    }
}