import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { usersAPI, achievementsAPI } from '../../api/api';
import './UserDetail.css';

const UserDetail = () => {
  const { uuid } = useParams();
  const navigate = useNavigate();
  const [user, setUser] = useState(null);
  const [userAchievements, setUserAchievements] = useState([]);
  const [allAchievements, setAllAchievements] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        setLoading(true);
        setError(null);
        
        const userResponse = await usersAPI.getUser(uuid);
        setUser(userResponse.data);
        
  
        const achievementsResponse = await usersAPI.getUserAchievements(uuid);
        setUserAchievements(achievementsResponse.data.achievements || []);
        
 
        const allAchievementsResponse = await achievementsAPI.getList();
        setAllAchievements(allAchievementsResponse.data.achievements || []);
      } catch (err) {
        console.error('Ошибка загрузки данных пользователя:', err);
        setError('Ошибка загрузки данных пользователя');
      } finally {
        setLoading(false);
      }
    };

    fetchUserData();
  }, [uuid]);


  const renderAchievements = () => {
    if (!userAchievements || userAchievements.length === 0) {
      return <p>У пользователя нет достижений</p>;
    }

    return (
      <div className="achievements-list">
        {userAchievements.map((achievement) => (
          <div 
            key={achievement.id} 
            className={`achievement-card ${achievement.achieved ? 'achieved' : 'not-achieved'}`}
          >
            <h4>{achievement.title}</h4>
            <p>{achievement.description}</p>
            <div className="achievement-progress">
              <div className="progress-bar">
                <div 
                  className="progress-fill" 
                  style={{ 
                    width: `${Math.min(100, (achievement.current_count / (achievement.achieved ? 1 : allAchievements.find(a => a.id === achievement.id)?.condition?.count || 1)) * 100)}%` 
                  }}
                ></div>
              </div>
              <span className="progress-text">
                {achievement.current_count} / 
                {achievement.achieved ? achievement.current_count : allAchievements.find(a => a.id === achievement.id)?.condition?.count || '?'}
              </span>
            </div>
            {achievement.achieved && (
              <p className="achievement-date">Получено: {new Date(achievement.achieved_at).toLocaleString()}</p>
            )}
          </div>
        ))}
      </div>
    );
  };

  if (loading) return <div className="user-detail"><p>Загрузка данных пользователя...</p></div>;
  if (error) return <div className="user-detail"><p className="error">{error}</p></div>;
  if (!user) return <div className="user-detail"><p>Пользователь не найден</p></div>;

  return (
    <div className="user-detail">
      <button onClick={() => navigate('/users')} className="back-button">
        ← Назад к списку пользователей
      </button>
      
      <div className="user-profile">
        <div className="user-header">
          <img 
            src={user.avatar ? user.avatar.trim() : '/default-avatar.png'} 
            alt={`${user.name} ${user.last_name}`} 
            className="user-avatar"
          />
          <div className="user-title">
            <h1>{user.name} {user.second_name} {user.last_name}</h1>
            <p className="user-login">@{user.login}</p>
          </div>
        </div>
        
        <div className="user-stats">
          <div className="stat-item">
            <span className="stat-label">Ранг</span>
            <span className="stat-value">{user.rank?.name || 'Не указан'}</span>
          </div>
          <div className="stat-item">
            <span className="stat-label">Очки</span>
            <span className="stat-value">{user.total_points || 0}</span>
          </div>
          <div className="stat-item">
            <span className="stat-label">Завершенные курсы</span>
            <span className="stat-value">{user.finished_courses || 0}</span>
          </div>
          <div className="stat-item">
            <span className="stat-label">Создан</span>
            <span className="stat-value">{new Date(user.created_at).toLocaleString()}</span>
          </div>
          <div className="stat-item">
            <span className="stat-label">Обновлен</span>
            <span className="stat-value">{new Date(user.updated_at).toLocaleString()}</span>
          </div>
        </div>
      </div>

      <div className="achievements-section">
        <h2>Достижения пользователя</h2>
        {renderAchievements()}
      </div>
    </div>
  );
};

export default UserDetail;