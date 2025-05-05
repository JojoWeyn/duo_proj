import { useNavigate } from 'react-router-dom';

const LogoutButton = () => {
  const navigate = useNavigate();

  const handleLogout = () => {
    // Удаляем токен из localStorage
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');


    // Перенаправляем пользователя на страницу логина
    navigate('/login');
  };

  return (
    <button onClick={handleLogout} className="logout-button">
      Выйти
    </button>
  );
};

export default LogoutButton;
