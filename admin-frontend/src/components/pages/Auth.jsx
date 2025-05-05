import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { authAPI } from '../../api/api';

const Auth = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [showResetModal, setShowResetModal] = useState(false);
  const [resetEmail, setResetEmail] = useState('');
  const [resetStep, setResetStep] = useState(1);
  const [resetCode, setResetCode] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [resetMessage, setResetMessage] = useState({ text: '', type: '' });
  const [resetLoading, setResetLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const response = await authAPI.login(email, password);
      localStorage.setItem('access_token', response.data.access_token);
      localStorage.setItem('refresh_token', response.data.refresh_token);
      navigate('/courses');
    } catch (err) {
      setError(err.response?.data?.message || 'Login failed. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleRequestCode = async (e) => {
    e.preventDefault();
    setResetLoading(true);
    setResetMessage({ text: '', type: '' });

    try {
      await authAPI.getVerificationCode(resetEmail);
      setResetMessage({ text: 'Код подтверждения отправлен на вашу почту', type: 'success' });
      setResetStep(2);
    } catch (error) {
      setResetMessage({ 
        text: error.response?.data?.message || 'Ошибка при отправке кода подтверждения', 
        type: 'error' 
      });
    } finally {
      setResetLoading(false);
    }
  };

  const handleResetPassword = async (e) => {
    e.preventDefault();
    
    if (newPassword !== confirmPassword) {
      setResetMessage({ text: 'Пароли не совпадают', type: 'error' });
      return;
    }

    if (newPassword.length < 8) {
      setResetMessage({ text: 'Пароль должен содержать не менее 8 символов', type: 'error' });
      return;
    }

    setResetLoading(true);
    setResetMessage({ text: '', type: '' });

    try {
      await authAPI.resetPassword(resetCode, resetEmail, newPassword);
      setResetMessage({ text: 'Пароль успешно изменен', type: 'success' });
      setTimeout(() => {
        setShowResetModal(false);
        setResetStep(1);
        setResetCode('');
        setNewPassword('');
        setConfirmPassword('');
        setResetMessage({ text: '', type: '' });
      }, 2000);
    } catch (error) {
      setResetMessage({ 
        text: error.response?.data?.message || 'Ошибка при сбросе пароля', 
        type: 'error' 
      });
    } finally {
      setResetLoading(false);
    }
  };

  const closeResetModal = () => {
    setShowResetModal(false);
    setResetStep(1);
    setResetEmail('');
    setResetCode('');
    setNewPassword('');
    setConfirmPassword('');
    setResetMessage({ text: '', type: '' });
  };

  const handleBack = () => {
    setResetStep(resetStep - 1);
    setResetMessage({ text: '', type: '' });
  };

  return (
    <div className="login-container">
        <h1>Admin Panel Login</h1>
        
        {error && (
          <div className="error-message">
            {error}
          </div>
        )}
         
         
        <form onSubmit={handleSubmit} className="admin-form">
            <div className="form-group">
                <label>Email:</label>
                    <input
                        type="email"
                        placeholder="Email Address"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        required
                    />
            </div>
            <div className="form-group">
                <label>Password:</label>
                    <input
                        type="password"
                        placeholder="Password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                    />
          </div>
            
          <button type="submit" className="submit-button" disabled={loading}>
            {loading ? 'Загрузка...' : 'Войти'}
          </button>
          
          <div className="reset-password-link">
            <button 
              type="button" 
              className="text-button" 
              onClick={() => setShowResetModal(true)}
            >
              Забыли пароль?
            </button>
          </div>
        </form>

        {showResetModal && (
          <div className="modal-overlay">
            <div className="modal-content">
              <button className="close-button" onClick={closeResetModal}>&times;</button>
              <h2>Сброс пароля</h2>
              
              {resetMessage.text && (
                <div className={`message ${resetMessage.type}`}>
                  {resetMessage.text}
                </div>
              )}

              {resetStep === 1 && (
                <form onSubmit={handleRequestCode} className="reset-form">
                  <div className="form-group">
                    <label>Email:</label>
                    <input
                      type="email"
                      placeholder="Введите ваш email"
                      value={resetEmail}
                      onChange={(e) => setResetEmail(e.target.value)}
                      required
                    />
                  </div>
                  <button 
                    type="submit" 
                    className="submit-button" 
                    disabled={resetLoading}
                  >
                    {resetLoading ? 'Отправка...' : 'Отправить код'}
                  </button>
                </form>
              )}

              {resetStep === 2 && (
                <form onSubmit={handleResetPassword} className="reset-form">
                  <div className="form-group">
                    <label>Код подтверждения:</label>
                    <input
                      type="text"
                      placeholder="Введите код из письма"
                      value={resetCode}
                      onChange={(e) => setResetCode(e.target.value)}
                      required
                    />
                  </div>
                  <div className="form-group">
                    <label>Новый пароль:</label>
                    <input
                      type="password"
                      placeholder="Введите новый пароль"
                      value={newPassword}
                      onChange={(e) => setNewPassword(e.target.value)}
                      required
                    />
                  </div>
                  <div className="form-group">
                    <label>Подтверждение пароля:</label>
                    <input
                      type="password"
                      placeholder="Повторите новый пароль"
                      value={confirmPassword}
                      onChange={(e) => setConfirmPassword(e.target.value)}
                      required
                    />
                  </div>
                  <div className="button-group">
                    <button 
                      type="button" 
                      className="back-button" 
                      onClick={handleBack}
                    >
                      Назад
                    </button>
                    <button 
                      type="submit" 
                      className="submit-button" 
                      disabled={resetLoading}
                    >
                      {resetLoading ? 'Сохранение...' : 'Сохранить новый пароль'}
                    </button>
                  </div>
                </form>
              )}
            </div>
          </div>
        )}
      </div>
    
  );
};

export default Auth;