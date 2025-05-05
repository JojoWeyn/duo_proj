import React, { useState, useEffect } from 'react';
import { authAPI } from '../../api/api';

const Security = () => {
  const [email, setEmail] = useState('');
  const [code, setCode] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [step, setStep] = useState(1);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState({ text: '', type: '' });
  const [userEmail, setUserEmail] = useState('');

  useEffect(() => {
    const fetchUserInfo = async () => {
      try {
        const response = await authAPI.getUserInfo();
        if (response.data && response.data.email) {
          setUserEmail(response.data.email);
          setEmail(response.data.email);
        }
      } catch (error) {
        console.error('Ошибка при получении информации о пользователе:', error);
      }
    };

    fetchUserInfo();
  }, []);

  // Запрос кода подтверждения
  const handleRequestCode = async (e) => {
    e.preventDefault();
    setLoading(true);
    setMessage({ text: '', type: '' });

    try {
      await authAPI.getVerificationCode(email);
      setMessage({ text: 'Код подтверждения отправлен на вашу почту', type: 'success' });
      setStep(2);
    } catch (error) {
      setMessage({ 
        text: error.response?.data?.message || 'Ошибка при отправке кода подтверждения', 
        type: 'error' 
      });
    } finally {
      setLoading(false);
    }
  };

  // Сброс пароля
  const handleResetPassword = async (e) => {
    e.preventDefault();
    
    if (newPassword !== confirmPassword) {
      setMessage({ text: 'Пароли не совпадают', type: 'error' });
      return;
    }

    if (newPassword.length < 8) {
      setMessage({ text: 'Пароль должен содержать не менее 8 символов', type: 'error' });
      return;
    }

    setLoading(true);
    setMessage({ text: '', type: '' });

    try {
      await authAPI.resetPassword(code, email, newPassword);
      setMessage({ text: 'Пароль успешно изменен', type: 'success' });
      setStep(1);
      setCode('');
      setNewPassword('');
      setConfirmPassword('');
    } catch (error) {
      setMessage({ 
        text: error.response?.data?.message || 'Ошибка при сбросе пароля', 
        type: 'error' 
      });
    } finally {
      setLoading(false);
    }
  };

  // Возврат к предыдущему шагу
  const handleBack = () => {
    setStep(step - 1);
    setMessage({ text: '', type: '' });
  };

  return (
    <div className="security-container">
      <h1>Безопасность</h1>
      <div className="security-card">
        <h2>Изменение пароля</h2>

        {message.text && (
          <div className={`message ${message.type}`}>
            {message.text}
          </div>
        )}

        {step === 1 && (
          <form onSubmit={handleRequestCode} className="security-form">
            <div className="form-group">
              <label>Email:</label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                disabled={userEmail !== ''}
                required
              />
            </div>
            <button type="submit" className="submit-button" disabled={loading}>
              {loading ? 'Отправка...' : 'Получить код подтверждения'}
            </button>
          </form>
        )}

        {step === 2 && (
          <form onSubmit={(e) => { e.preventDefault(); setStep(3); }} className="security-form">
            <div className="form-group">
              <label>Код подтверждения:</label>
              <input
                type="text"
                value={code}
                onChange={(e) => setCode(e.target.value)}
                required
              />
            </div>
            <div className="button-group">
              <button type="button" className="back-button" onClick={handleBack}>
                Назад
              </button>
              <button type="submit" className="submit-button">
                Далее
              </button>
            </div>
          </form>
        )}

        {step === 3 && (
          <form onSubmit={handleResetPassword} className="security-form">
            <div className="form-group">
              <label>Новый пароль:</label>
              <input
                type="password"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                required
              />
            </div>
            <div className="form-group">
              <label>Подтвердите пароль:</label>
              <input
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                required
              />
            </div>
            <div className="button-group">
              <button type="button" className="back-button" onClick={handleBack}>
                Назад
              </button>
              <button type="submit" className="submit-button" disabled={loading}>
                {loading ? 'Сохранение...' : 'Изменить пароль'}
              </button>
            </div>
          </form>
        )}
      </div>

      <style jsx>{`
        .security-container {
          max-width: 600px;
          margin: 0 auto;
          padding: 20px;
        }
        
        .security-card {
          background: white;
          border-radius: 8px;
          padding: 20px;
          box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
        }
        
        .security-form {
          display: flex;
          flex-direction: column;
          gap: 15px;
        }
        
        .form-group {
          display: flex;
          flex-direction: column;
          gap: 5px;
        }
        
        .form-group label {
          font-weight: 500;
        }
        
        .form-group input {
          padding: 10px;
          border: 1px solid #ddd;
          border-radius: 4px;
          font-size: 16px;
        }
        
        .button-group {
          display: flex;
          gap: 10px;
          margin-top: 10px;
        }
        
        .submit-button {
          background-color: #4a90e2;
          color: white;
          border: none;
          padding: 10px 15px;
          border-radius: 4px;
          cursor: pointer;
          font-size: 16px;
        }
        
        .submit-button:disabled {
          background-color: #cccccc;
          cursor: not-allowed;
        }
        
        .back-button {
          background-color: #f5f5f5;
          color: #333;
          border: 1px solid #ddd;
          padding: 10px 15px;
          border-radius: 4px;
          cursor: pointer;
          font-size: 16px;
        }
        
        .message {
          padding: 10px;
          border-radius: 4px;
          margin-bottom: 15px;
        }
        
        .success {
          background-color: #e7f7e7;
          color: #2e7d32;
          border: 1px solid #c8e6c9;
        }
        
        .error {
          background-color: #ffebee;
          color: #c62828;
          border: 1px solid #ffcdd2;
        }
      `}</style>
    </div>
  );
};

export default Security;