import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { achievementsAPI } from '../../api/api';

export const AchievementCreate = () => {
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    action: 'lesson',
    count: 1,
    secret: false
  });
  const [errors, setErrors] = useState({});
  const navigate = useNavigate();

  const validateForm = () => {
    const newErrors = {};
    if (!formData.title) newErrors.title = 'Название обязательно';
    if (!formData.description) newErrors.description = 'Описание обязательно';
    if (formData.count < 1) newErrors.count = 'Минимальное значение 1';
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!validateForm()) return;

    try {
      await achievementsAPI.create({
        title: formData.title,
        description: formData.description,
        condition: {
          action: formData.action,
          count: formData.count
        },
        secret: formData.secret
      });
      navigate('/achievements');
    } catch (error) {
      console.error('Ошибка создания:', error);
      setErrors({ submit: 'Ошибка при создании ачивки' });
    }
  };

  return (
    <div className="form-container">
      <h2>Создание новой ачивки</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Название:</label>
          <input
            type="text"
            value={formData.title}
            onChange={(e) => setFormData({...formData, title: e.target.value})}
          />
          {errors.title && <span className="error">{errors.title}</span>}
        </div>

        <div className="form-group">
          <label>Описание:</label>
          <textarea
            value={formData.description}
            onChange={(e) => setFormData({...formData, description: e.target.value})}
          />
          {errors.description && <span className="error">{errors.description}</span>}
        </div>

        <div className="form-group">
          <label>Тип действия:</label>
          <select
            value={formData.action}
            onChange={(e) => setFormData({...formData, action: e.target.value})}
          >
            <option value="lesson">Прохождение Урока</option>
            <option value="exercise">Прохождение Упражнения</option>
            <option value="update">Изменение профиля</option>
            <option value="course">Прохождение Курса</option>
            <option value="question">Прохождение Вопроса</option>
            <option value="login">Вхождение в аккаунт</option>
          </select>
        </div>

        <div className="form-group">
          <label>Количество:</label>
          <input
            type="number"
            min="1"
            value={formData.count}
            onChange={(e) => setFormData({...formData, count: parseInt(e.target.value)})}
          />
          {errors.count && <span className="error">{errors.count}</span>}
        </div>

        <div className="form-group checkbox">
          <label>
            <input
              type="checkbox"
              checked={formData.secret}
              onChange={(e) => setFormData({...formData, secret: e.target.checked})}
            />
            Секретная ачивка
          </label>
        </div>

        {errors.submit && <div className="error-message">{errors.submit}</div>}

        <div className="form-actions">
          <button type="button" onClick={() => navigate('/achievements')} className="cancel-btn">
            Отмена
          </button>
          <button type="submit" className="submit-btn">
            Создать
          </button>
        </div>
      </form>
    </div>
  );
};