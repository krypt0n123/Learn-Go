// frontend/src/main.ts (新内容)

import * as api from './api.js'; // 导入所有 API 函数
import * as ui from './ui.js'; // 导入所有 UI 函数和元素

/**
 * 主刷新函数：获取最新数据并重新渲染 UI
 */
async function refreshTodoList() {
    try {
        // 1. 从 API 获取数据
        const todos = await api.getAllTodos();
        // 2. 使用 UI 模块渲染
        ui.renderTodos(todos);
    } catch (error) {
        console.error('获取 Todos 时出错:', error);
        const message = error instanceof Error ? error.message : '加载失败';
        ui.renderError(message);
    }
}

/**
 * 处理表单提交
 */
async function handleFormSubmit(event: SubmitEvent) {
    event.preventDefault();
    if (!ui.todoNameInput || !ui.todoNameInput.value) {
        alert('请输入待办事项名称');
        return;
    }

    // V V V V V 修改开始 V V V V V

    // 1. 从两个输入框读取数据
    const newTodoName = ui.todoNameInput.value;
    const newTodoDescription = ui.todoDescriptionInput ? ui.todoDescriptionInput.value : "";

    try {
        // 2. 将两个数据都发送到 API
        await api.createTodo(newTodoName, newTodoDescription);

        // 3. 清空两个输入框
        ui.todoNameInput.value = '';
        if (ui.todoDescriptionInput) {
            ui.todoDescriptionInput.value = '';
        }

        // 4. 刷新列表
        await refreshTodoList();
        // ^ ^ ^ ^ ^ 修改结束 ^ ^ ^ ^ ^

    } catch (error) {
        console.error('创建 Todo 时出错:', error);
        alert('创建失败，请查看控制台');
    }
}

/**
 * 处理列表中的点击事件 (删除 和 切换)
 */
async function handleListClick(event: MouseEvent) {
    const target = event.target as HTMLElement;
    const action = target.dataset.action;
    if (!action) return; // 点击的不是 action 按钮

    const li = target.closest('li');
    if (!li || !li.dataset.id) return; // 找不到父元素 <li> 或 ID

    const todoId = li.dataset.id;

    try {
        if (action === 'delete') {
            // 1. 调用 API 删除
            await api.deleteTodo(todoId);
            // 2. 刷新
            await refreshTodoList();
        }

        if (action === 'toggle') {
            // 从 data 属性读取数据
            const name = li.dataset.name!;
            const description = li.dataset.description || "";
            const checkbox = target as HTMLInputElement;
            const newCompletedStatus = checkbox.checked; // 获取*新*的状态

            // 1. 调用 API 更新
            await api.updateTodo(todoId, name, description, newCompletedStatus);
            // 2. 刷新
            await refreshTodoList();
        }
    } catch (error) {
        console.error("操作失败:", error);
        alert("操作失败，请查看控制台");
    }
}

// ---------------------------------
// 程序主入口
// ---------------------------------
document.addEventListener('DOMContentLoaded', () => {
    // 1. 页面加载时，立即刷新一次列表
    refreshTodoList();

    // 2. 监听表单提交
    if (ui.newTodoForm) {
        ui.newTodoForm.addEventListener('submit', handleFormSubmit);
    }

    // 3. 监听整个列表的点击事件
    if (ui.todoListElement) {
        ui.todoListElement.addEventListener('click', handleListClick);
    }
});