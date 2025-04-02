import { useState } from 'react';
import { questionsAPI } from '../../api/api';
import { useParams, useNavigate } from 'react-router-dom';

export const QuestionOptionCreate = () => {
  const { questionUUID } = useParams();  // Получаем UUID вопроса из параметров URL
  const navigate = useNavigate();
  
  const [text, setText] = useState('');
  const [isCorrect, setIsCorrect] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    const optionData = {
      text,
      is_correct: isCorrect,
      question_uuid: questionUUID 
    };

    try {
      await questionsAPI.createQuestionOption(optionData);
      navigate(`/questions/${questionUUID}`); 
    } catch (err) {
      setError('Ошибка при создании варианта ответа!');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="question-option-create">
      <h2>Создать вариант ответа</h2>
      <form onSubmit={handleSubmit}>
        {error && <p className="error-message">{error}</p>}
        
        <div className="form-group">
          <label>Текст варианта ответа</label>
          <input
            type="text"
            value={text}
            onChange={(e) => setText(e.target.value)}
            required
          />
        </div>
        
        <div className="form-group">
          <label>Правильный вариант</label>
          <input
            type="checkbox"
            checked={isCorrect}
            onChange={(e) => setIsCorrect(e.target.checked)}
          />
        </div>

        <button type="submit" disabled={loading}>
          {loading ? 'Создание варианта...' : 'Создать вариант ответа'}
        </button>
      </form>
    </div>
  );
};
