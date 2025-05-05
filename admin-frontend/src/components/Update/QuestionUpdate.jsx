import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { questionsAPI } from '../../api/api';

export const QuestionUpdate = () => {
  const { uuid } = useParams();
  const navigate = useNavigate();
  const [question, setQuestion] = useState({
    text: '',
    type_id: 1,
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchQuestion = async () => {
      try {
        const response = await questionsAPI.getQuestionMeta(uuid);
        setQuestion({
          text: response.data.text,
          type_id: response.data.type.id,
        });
        setLoading(false);
      } catch (err) {
        setError('Не удалось загрузить вопрос');
        setLoading(false);
      }
    };
    fetchQuestion();
  }, [uuid]);


  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await questionsAPI.updateQuestion(uuid, {
        uuid: uuid,
        text: question.text,
        type_id: question.type_id,
      });
      navigate(-1);
    } catch (err) {
      setError('Ошибка при обновлении вопроса');
    }
  };

  if (loading) return <p>Загрузка...</p>;
  if (error) return <p className="error-message">{error}</p>;

  return (
    <div className="admin-panel">
      <h2>Редактировать вопрос</h2>
      <form className="admin-form" onSubmit={handleSubmit}>
        {error && <p className="error-message">{error}</p>}
        <div className="form-group">
          <label>Текст вопроса</label>
          <textarea
            value={question.text}
            onChange={(e) => setQuestion({...question, text: e.target.value})}
            required
          ></textarea>
        </div>
        <div className="form-group">
          <label>Тип вопроса</label>
          <select
            value={question.type_id}
            onChange={(e) => setQuestion({...question, type_id: Number(e.target.value)})}
          >
            <option value={1}>Одиночный выбор</option>
            <option value={2}>Множественный выбор</option>
            <option value={3}>Сопоставление</option>
          </select>
        </div>
        <button className="submit-button" type="submit" disabled={loading}>
          Обновить вопрос
        </button>
      </form>
    </div>
  );
};
