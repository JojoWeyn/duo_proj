import { Navigate, Outlet } from 'react-router-dom';
import { useState, useEffect } from 'react';
import { authAPI } from '../api/api';

const PrivateRoute = () => {
  const [isLoading, setIsLoading] = useState(true);
  const [isValid, setIsValid] = useState(false);

  useEffect(() => {
    const validateToken = async () => {
      const accessToken = localStorage.getItem('access_token');
      
      if (!accessToken) {
        setIsValid(false);
        setIsLoading(false);
        return;
      }

      try {
        await authAPI.checkToken();
        setIsValid(true);
      } catch (error) {
        console.error('Ошибка проверки токена:', error);
        localStorage.removeItem('access_token');
        setIsValid(false);
      } finally {
        setIsLoading(false);
      }
    };

    validateToken();
  }, []);

  if (isLoading) {
    return <div>Загрузка...</div>;
  }

  if (!isValid) {
    return <Navigate to="/login" replace />;
  }

  return <Outlet />;
};

export default PrivateRoute;
