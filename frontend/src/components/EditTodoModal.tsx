import React, { useState, useEffect, type FormEvent } from 'react';
import type { Todo, UpdateTodoRequest } from '../types';

interface EditTodoModalProps {
  todo: Todo | null;
  onClose: () => void;
  onSave: (id: number, data: UpdateTodoRequest) => Promise<void>;
}

export const EditTodoModal: React.FC<EditTodoModalProps> = ({
  todo,
  onClose,
  onSave,
}) => {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    if (todo) {
      setTitle(todo.title);
      setDescription(todo.description || '');
    }
  }, [todo]);

  if (!todo) {
    return null;
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!title.trim() || isSaving) return;

    setIsSaving(true);
    try {
      await onSave(todo.id, {
        title: title.trim(),
        description: description.trim(),
      });
      onClose();
    } catch (error) {
      console.error('Ошибка при сохранении задачи:', error);
    } finally {
      setIsSaving(false);
    }
  };

  const isDirty =
    title.trim() !== todo.title.trim() ||
    description.trim() !== (todo.description || '').trim();

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <h2>Редактировать задачу</h2>
        <form className="modal-form" onSubmit={handleSubmit}>
          <label>
            Заголовок:
            <input
              type="text"
              className="form-input"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
            />
          </label>
          <label>
            Описание:
            <textarea
              className="form-input"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={4}
            />
          </label>
          <div className="modal-actions">
            <button
              type="button"
              className="btn btn-secondary"
              onClick={onClose}
              disabled={isSaving}
            >
              Отмена
            </button>
            <button
              type="submit"
              className="btn btn-primary"
              disabled={!isDirty || isSaving || !title.trim()}
            >
              {isSaving ? 'Сохранение...' : 'Сохранить'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};