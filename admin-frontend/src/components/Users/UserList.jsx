import React, { useState, useEffect } from 'react';
import { usersAPI } from '../../api/api';


const UserList = () => {
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
  
    // Функция для рендеринга каждого пользователя
    const renderUser = (user) => (
      <div key={user.uuid} className="user-card">
        <img src={user.avatar || ''} alt="Avatar"/>
        <div className="user-info">
          <h3>{user.name} {user.second_name} {user.last_name}</h3>
          <p>Rank: {user.rank.name}</p>
          <p>Created at: {new Date(user.created_at).toLocaleString()}</p>
          <p>Updated at: {new Date(user.updated_at).toLocaleString()}</p>
        </div>
      </div>
    );
  
    return (
      <div className="user-list">
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