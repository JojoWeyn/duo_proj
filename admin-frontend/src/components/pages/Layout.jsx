import React, { useState, useEffect } from 'react';
import LogoutButton from '../LogoutButton';
import { Outlet, useNavigate } from 'react-router-dom';
import './Layout.css';

const Layout = () => {
  const navigate = useNavigate();
  const [menuOpen, setMenuOpen] = useState(false);

  const toggleMenu = () => {
    setMenuOpen(!menuOpen);
  };

  // Закрываем меню при изменении размера экрана на десктоп
  useEffect(() => {
    const handleResize = () => {
      if (window.innerWidth > 768 && menuOpen) {
        setMenuOpen(false);
      }
    };

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, [menuOpen]);

  return (
    <div className="layout-container">
      <button className="mobile-menu-button" onClick={toggleMenu}>
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          {menuOpen ? (
            <>
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </>
          ) : (
            <>
              <line x1="3" y1="12" x2="21" y2="12"></line>
              <line x1="3" y1="6" x2="21" y2="6"></line>
              <line x1="3" y1="18" x2="21" y2="18"></line>
            </>
          )}
        </svg>
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
