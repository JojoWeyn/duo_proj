import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { achievementsAPI } from '../../api/api';

export const AchievementUpdate = () => {
  const { id } = useParams();
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    action: 'lesson',
    count: 1,
    secret: false
  });
  const [errors, setErrors] = useState({});
  const navigate = useNavigate();

  useEffect(() => {
    const fetchAchievement = async () => {
      try {
        const response = await achievementsAPI.get(id);
        const ach = response.data.achievement;
        setFormData({
          title: ach.title,
          description: ach.description,
          action: ach.condition.action,
          count: ach.condition.count,
          secret: ach.secret
        });
      } catch (error) {
        console.error('Ошибка загрузки:', error);
        navigate('/achievements');
      }
    };
    fetchAchievement();
  }, [id]);

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
      await achievementsAPI.update(id, {
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
      console.error('Ошибка обновления:', error);
      setErrors({ submit: 'Ошибка при обновлении ачивки' });
    }
  };

  return (
    <div className="form-container">
      <h2>Редактирование ачивки</h2>
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
            <option value="lesson">Урок</option>
            <option value="test">Тест</option>
            <option value="course">Курс</option>
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
            Сохранить
          </button>
        </div>
      </form>
    </div>
  );
};