import { Todo } from './types.js'; // 导入我们的类型

// 获取所有 todos
export async function getAllTodos(): Promise<Todo[]> {
    const response = await fetch('/get-all-todos');
    if (!response.ok) {
        throw new Error('获取数据失败');
    }
    return response.json();
}

// 创建一个 new todo
export async function createTodo(name: string, description: string): Promise<void> {
    const response = await fetch('/create', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            name: name,
            description: description
        })
    });
    if (!response.ok) {
        throw new Error('创建失败');
    }
}

// 更新一个 todo
export async function updateTodo(id: string, name: string, description: string, completed: boolean): Promise<void> {
    const response = await fetch('/update', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            id: id,
            name: name,
            description: description,
            completed: String(completed) // 你的 Go 后端期望的是字符串
        })
    });
    if (!response.ok) {
        throw new Error('更新失败');
    }
}

// 删除一个 todo
export async function deleteTodo(id: string): Promise<void> {
    const response = await fetch('/delete', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({ id: id })
    });
    if (!response.ok) {
        throw new Error('删除失败');
    }
}