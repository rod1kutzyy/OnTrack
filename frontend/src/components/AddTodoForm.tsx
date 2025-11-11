import React, { useState, type FormEvent } from 'react';

interface AddTodoFormProps {
  onAddTodo: (title: string, description?: string) => Promise<void>;
}

export const AddTodoForm: React.FC<AddTodoFormProps> = ({ onAddTodo }) => {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;

    setIsSubmitting(true);
    await onAddTodo(title.trim(), description.trim() || undefined);
    
    setIsSubmitting(false);
    setTitle('');
    setDescription('');
  };

  return (
    <form className="add-todo-form" onSubmit={handleSubmit}>
      <input
        type="text"
        className="form-input"
        placeholder="Новая задача..."
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        required
      />
      <input
        type="text"
        className="form-input"
        placeholder="Описание (необязательно)"
        value={description}
        onChange={(e) => setDescription(e.target.value)}
      />
      <button
        type="submit"
        className="btn btn-primary"
        disabled={isSubmitting || !title.trim()}
      >
        {isSubmitting ? '...' : 'Добавить'}
      </button>
    </form>
  );
};