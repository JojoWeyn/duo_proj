import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { lessonsAPI, coursesAPI } from '../../api/api';

export const LessonUpdate = () => {
  const { lessonUUID } = useParams();  // Получаем UUID урока из параметров маршрута
  const [lesson, setLesson] = useState({
    title: '',
    description: '',
    difficulty_id: 1,
    order: 1,
    course_uuid: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  // Функция для загрузки данных урока
  const fetchLesson = async () => {
    try {
      const response = await lessonsAPI.getLesson(lessonUUID);
      setLesson(response.data)  
      console.log(response.data)
    } catch (err) {
      console.error('Ошибка при загрузке урока:', err);
      setError('Ошибка при загрузке данных урока');
    }
  };

  useEffect(() => {
    fetchLesson();  // Загружаем данные при монтировании компонента
  }, [lessonUUID]);

  // Обработчик изменения данных в форме
  const handleChange = (e) => {
    const { name, value } = e.target;
    setLesson(prev => ({
      ...prev,
      [name]: value,
    }));
  };

  // Обработчик отправки формы
  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
        const payload = {
            uuid: lessonUUID,
            title: lesson.title,
            description: lesson.description,
            difficulty_id: parseInt(lesson.difficulty_id),
            order: parseInt(lesson.order),
            course_uuid: lesson.course_uuid
          };
      await lessonsAPI.updateLesson(lessonUUID, payload);  // Отправка обновленных данных на сервер
      navigate(`/courses/${lesson.course_uuid}`);  // Перенаправление на страницу урока
    } catch (err) {
      console.error('Ошибка при обновлении урока:', err);
      setError('Ошибка при обновлении урока');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="admin-panel">
      <h2>Редактирование урока</h2>

      {error && <p className="error-message">{error}</p>}

      <form className="admin-form" onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="title">Название урока</label>
          <input
            type="text"
            id="title"
            name="title"
            value={lesson.title}
            onChange={handleChange}
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="description">Описание</label>
          <textarea
            id="description"
            name="description"
            value={lesson.description}
            onChange={handleChange}
            rows={8}
            required
          ></textarea>
        </div>

        <div className="form-group">
          <label htmlFor="difficulty_id">Сложность</label>
          <select
            id="difficulty_id"
            name="difficulty_id"
            value={lesson.difficulty_id}
            onChange={handleChange}
            required
          >
            <option value={1}>Легкий</option>
            <option value={2}>Средний</option>
            <option value={3}>Сложный</option>
          </select>
        </div>

        <div className="form-group">
          <label htmlFor="order">Порядковый номер</label>
          <input
            type="number"
            id="order"
            name="order"
            value={lesson.order}
            onChange={handleChange}
            required
          />
        </div>

        <button className='submit-button' type="submit" disabled={loading}>
          {loading ? 'Обновление...' : 'Обновить урок'}
        </button>
      </form>
    </div>
  );
};
