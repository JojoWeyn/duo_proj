import { Navigate, Outlet } from 'react-router-dom';

const PrivateRoute = () => {
  const accessToken = localStorage.getItem('access_token');

  // Если токен отсутствует, перенаправляем на страницу логина
  if (!accessToken) {
    return <Navigate to="/login" replace />;
  }

  // Если токен есть, отображаем дочерние маршруты (компоненты)
  return <Outlet />;
};

export default PrivateRoute;
