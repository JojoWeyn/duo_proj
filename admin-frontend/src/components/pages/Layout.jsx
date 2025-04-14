import React, { useState } from 'react';
import LogoutButton from '../LogoutButton';
import { Outlet, useNavigate } from 'react-router-dom';

const Layout = () => {
  const navigate = useNavigate();
  const [menuOpen, setMenuOpen] = useState(false);

  const toggleMenu = () => {
    setMenuOpen(!menuOpen);
  };

  return (
    <div className="layout-container">
      <button className="mobile-menu-button" onClick={toggleMenu}>
        
      </button>

      <div className={`dashboard-buttons ${menuOpen ? 'open' : ''}`}>
        <h2>Z-School Admin</h2>
        <button className='dashboard-button' onClick={() => navigate('/courses')}>Курсы</button>
        <button className='dashboard-button' onClick={() => navigate('/users')}>Пользователи</button>
        <button className='dashboard-button' onClick={() => navigate('/files')}>Файлы</button>
        <LogoutButton />
      </div>

      <div className="layout-content">
        <Outlet />
      </div>
    </div>
  );
};

export default Layout;
