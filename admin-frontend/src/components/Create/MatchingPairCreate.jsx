import { useState } from 'react';
import { questionsAPI } from '../../api/api';
import { useParams, useNavigate } from 'react-router-dom';

export const MatchingPairCreate = () => {
  const { questionUUID } = useParams();
  const navigate = useNavigate();
  
  const [leftText, setLeftText] = useState('');
  const [rightText, setRightText] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    const matchingPairData = {
      left_text: leftText,
      right_text: rightText,
      question_uuid: questionUUID 
    };

    try {
      await questionsAPI.createMatchingPair(matchingPairData);
      navigate(`/questions/${questionUUID}`);  
    } catch (err) {
      setError('Ошибка при создании пары соответствия!');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="matching-pair-create">
      <h2>Создать пару соответствия</h2>
      <form onSubmit={handleSubmit}>
        {error && <p className="error-message">{error}</p>}
        
        <div className="form-group">
          <label>Текст слева</label>
          <input
            type="text"
            value={leftText}
            onChange={(e) => setLeftText(e.target.value)}
            required
          />
        </div>
        
        <div className="form-group">
          <label>Текст справа</label>
          <input
            type="text"
            value={rightText}
            onChange={(e) => setRightText(e.target.value)}
            required
          />
        </div>

        <button type="submit" disabled={loading}>
          {loading ? 'Создание пары...' : 'Создать пару соответствия'}
        </button>
      </form>
    </div>
  );
};
