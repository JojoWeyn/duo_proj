import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { achievementsAPI } from '../../api/api';

export const AchievementList = () => {
  const [achievements, setAchievements] = useState([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    fetchAchievements();
  }, []);

  const fetchAchievements = async () => {
    try {
      const response = await achievementsAPI.getList();
      setAchievements(response.data.achievements?.sort((a, b) => a.id - b.id) || []);
      setLoading(false);
    } catch (error) {
      console.error('Ошибка загрузки:', error);
      setLoading(false);
    }
  };

  const handleDelete = async (id) => {
    if (window.confirm('Удалить ачивку?')) {
      try {
        await achievementsAPI.delete(id);
        fetchAchievements();
      } catch (error) {
        console.error('Ошибка удаления:', error);
      }
    }
  };

  if (loading) return <div>Загрузка...</div>;

  return (
    <div className="table-container">
      <div className="table-header">
        <h2>Управление ачивками</h2>
      </div>

      <table className="data-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Название</th>
            <th>Описание</th>
            <th>Условие</th>
            <th>Секретная</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          {achievements.map((ach) => (
            <tr key={ach.id}>
              <td>{ach.id}</td>
              <td>{ach.title}</td>
              <td>{ach.description}</td>
              <td>
                <pre>{JSON.stringify(ach.condition, null, 2)}</pre>
              </td>
              <td>{ach.secret ? 'Да' : 'Нет'}</td>
              <td>
                <button
                  className="edit-btn"
                  onClick={() => navigate(`/achievements/${ach.id}/update`)}
                >
                  ✏️
                </button>
                <button
                  className="delete-btn"
                  onClick={() => handleDelete(ach.id)}
                >
                  🗑️
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <button 
          className="create-button"
          onClick={() => navigate('/achievements/create')}
        >
          + Новая ачивка
        </button>
    </div>
  );
};