import React from 'react';
import type { Todo } from '../types';

interface TodoItemProps {
  todo: Todo;
  onToggle: (id: number) => void;
  onDelete: (id: number) => void;
  onEdit: (todo: Todo) => void;
}

export const TodoItem: React.FC<TodoItemProps> = ({
  todo,
  onToggle,
  onDelete,
  onEdit,
}) => (
  <div className="todo-item" data-completed={todo.completed}>
    <input
      type="checkbox"
      className="todo-checkbox"
      checked={todo.completed}
      onChange={() => onToggle(todo.id)}
      title={
        todo.completed ? 'Отметить как невыполненное' : 'Отметить как выполненное'
      }
    />
    <div className="todo-content">
      <div className="todo-title">{todo.title}</div>
      {todo.description && (
        <div className="todo-description">{todo.description}</div>
      )}
    </div>
    {}
    <div className="todo-actions">
      <button
        className="btn-edit"
        onClick={() => onEdit(todo)}
        title="Редактировать"
      >
        ✏️
      </button>
      <button
        className="btn-danger"
        onClick={() => onDelete(todo.id)}
        title="Удалить"
      >
        &times;
      </button>
    </div>
  </div>
);