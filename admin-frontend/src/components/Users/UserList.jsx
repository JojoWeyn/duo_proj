import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { usersAPI } from '../../api/api';
import './UserList.css';


const UserList = () => {
    const navigate = useNavigate();
    const [users, setUsers] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
  
    // Функция для загрузки данных пользователей
    const loadUsers = async () => {
      try {
        setLoading(true);
        setError(null);
        const response = await usersAPI.getAll();
        setUsers(response.data.users); // предполагаем, что ответ содержит массив "users"
      } catch (err) {
        setError('Ошибка загрузки данных');
      } finally {
        setLoading(false);
      }
    };
  
    useEffect(() => {
      loadUsers();
    }, []);
  
    // Функция для перехода на страницу с подробной информацией о пользователе
    const handleUserClick = (uuid) => {
      navigate(`/users/${uuid}`);
    };

    // Функция для рендеринга каждого пользователя
    const renderUser = (user) => (
      <div 
        key={user.uuid} 
        className="user-card" 
        onClick={() => handleUserClick(user.uuid)}
        style={{ cursor: 'pointer' }}
      >
        <img src={user.avatar || '/default-avatar.png'} alt={`${user.name} ${user.last_name}`}/>
        <div className="user-info">
          <h3>{user.name} {user.second_name} {user.last_name}</h3>
          <p>Ранг: {user.rank.name}</p>
          <p>Создан: {new Date(user.created_at).toLocaleString()}</p>
          <p>Обновлен: {new Date(user.updated_at).toLocaleString()}</p>
        </div>
      </div>
    );
  
    return (
      <div className="user-list">
        <h1>Все пользователи</h1>
        {loading && <p>Загрузка...</p>}
        {error && <p>{error}</p>}
        {!loading && !error && users.length === 0 && <p>Нет пользователей</p>}
        <div className="user-list-container">
          {users.map(renderUser)}
        </div>
      </div>
    );
  };
  
  export default UserList;