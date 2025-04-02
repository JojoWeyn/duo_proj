import React from 'react';
import LogoutButton from '../../components/LogoutButton';
import { useParams, useNavigate } from 'react-router-dom';

const Dashboard = () => {
  const navigate = useNavigate();
  return (
    <div className="dashboard-container">
      <div className="dashboard-content">
        <h1 className="dashboard-title">Z-School Admin</h1>
        <div className="dashboard-buttons">
          <button className="dashboard-button" onClick={() => navigate('/courses')}>Курсы</button>
          <button className="dashboard-button" onClick={() => navigate('/users')}>Пользователи</button>
          <LogoutButton/>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
