import { useState, useEffect, useCallback } from 'react';

import type {
  Todo,
  CreateTodoRequest,
  UpdateTodoRequest,
  ApiSuccessResponse,
  GetTodosApiResponseData,
} from './types';

import { useTheme } from './hooks/useTheme';

import { ThemeToggle } from './components/ThemeToggle';
import { AddTodoForm } from './components/AddTodoForm';
import { TodoList } from './components/TodoList';
import { Pagination } from './components/Pagination';
import { EditTodoModal } from './components/EditTodoModal';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
const TODOS_PER_PAGE = 5;

function App() {
  const [theme, toggleTheme] = useTheme();
  const [todos, setTodos] = useState<Todo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);

  const [editingTodo, setEditingTodo] = useState<Todo | null>(null);


  const fetchTodos = useCallback(async (page: number) => {
    setLoading(true);
    setError(null);
    try {
      const params = new URLSearchParams({
        page: String(page),
        limit: String(TODOS_PER_PAGE),
      });

      const response = await fetch(`${API_BASE_URL}/todos?${params.toString()}`);
      if (!response.ok) {
        throw new Error(`Ошибка ${response.status}: ${response.statusText}`);
      }

      const result: ApiSuccessResponse<GetTodosApiResponseData> =
        await response.json();

      if (result.success && result.data) {
        setTodos(result.data.items);
        setTotalPages(result.data.pagination.total_pages);
        setCurrentPage(result.data.pagination.current_page);
      } else {
        throw new Error(
          result.message || 'API вернул ошибку в структуре данных'
        );
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Неизвестная ошибка');
      console.error('Ошибка при загрузке задач:', err);
    } finally {
      setLoading(false);
    }
  }, []);


  useEffect(() => {
    fetchTodos(currentPage);
  }, [currentPage, fetchTodos]);


  const handleAddTodo = async (title: string, description?: string) => {
    try {
      const newTodo: CreateTodoRequest = { title, description };
      const response = await fetch(`${API_BASE_URL}/todos`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newTodo),
      });
      if (!response.ok) throw new Error('Не удалось создать задачу');

      if (currentPage !== 1) {
        setCurrentPage(1);
      } else {
        await fetchTodos(1);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка при добавлении');
      console.error('Ошибка добавления:', err);
    }
  };


  const handleToggleTodo = async (id: number) => {
    setTodos((prevTodos) =>
      prevTodos.map((t) =>
        t.id === id ? { ...t, completed: !t.completed } : t
      )
    );
    try {
      const response = await fetch(`${API_BASE_URL}/todos/${id}/toggle`, {
        method: 'PATCH',
      });
      if (!response.ok) throw new Error('Не удалось обновить статус');

      await fetchTodos(currentPage);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка обновления');
      console.error('Ошибка обновления:', err);
    }
  };

  const handleDeleteTodo = async (id: number) => {
    const originalTodos = [...todos];
    setTodos((prevTodos) => prevTodos.filter((t) => t.id !== id));
    try {
      const response = await fetch(`${API_BASE_URL}/todos/${id}`, {
        method: 'DELETE',
      });
      if (!response.ok) {
        setTodos(originalTodos);
        throw new Error('Не удалось удалить задачу');
      }

      if (todos.length === 1 && currentPage > 1) {
        setCurrentPage(currentPage - 1);
      } else {
        await fetchTodos(currentPage);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка удаления');
      console.error('Ошибка удаления:', err);
    }
  };


  const handleUpdateTodo = async (
    id: number,
    data: UpdateTodoRequest
  ): Promise<void> => {
    try {
      const response = await fetch(`${API_BASE_URL}/todos/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error('Не удалось обновить задачу');
      }

      await fetchTodos(currentPage);
    } catch (err) {
      const errorMsg =
        err instanceof Error ? err.message : 'Ошибка при обновлении задачи';
      setError(errorMsg);
      console.error('Ошибка обновления:', err);
      throw err;
    }
  };


  const handleOpenEditModal = (todo: Todo) => {
    setEditingTodo(todo);
  };

  const handleCloseEditModal = () => {
    setEditingTodo(null);
  };

  const handlePageChange = (newPage: number) => {
    setCurrentPage(newPage);
  };


  return (
    <>
      <header className="app-header">
        <h1>OnTrack</h1>
        <ThemeToggle theme={theme} toggleTheme={toggleTheme} />
      </header>

      <main>
        <AddTodoForm onAddTodo={handleAddTodo} />

        {error && <p style={{ color: 'red' }}>Ошибка: {error}</p>}

        <TodoList
          todos={todos}
          loading={loading}
          onToggle={handleToggleTodo}
          onDelete={handleDeleteTodo}
          onEdit={handleOpenEditModal}
        />

        <Pagination
          currentPage={currentPage}
          totalPages={totalPages}
          onPageChange={handlePageChange}
        />
      </main>

      <EditTodoModal
        todo={editingTodo}
        onClose={handleCloseEditModal}
        onSave={handleUpdateTodo}
      />
    </>
  );
}

export default App;