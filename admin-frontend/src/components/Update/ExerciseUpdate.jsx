import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { exercisesAPI, lessonsAPI } from '../../api/api';

export const ExerciseUpdate = () => {
  const { exerciseUUID } = useParams();
  const navigate = useNavigate();

  const [exercise, setExercise] = useState({
    title: '',
    description: '',
    points: 0,
    order: 1,
    lesson_uuid: ''
  });

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchExercise = async () => {
      try {
        const res = await exercisesAPI.getExercise(exerciseUUID);
        setExercise(res.data);
      } catch (err) {
        console.error('Ошибка при загрузке упражнения:', err);
        setError('Не удалось загрузить данные упражнения');
      }
    };

    fetchExercise();
  }, [exerciseUUID]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setExercise(prev => ({
      ...prev,
      [name]: name === 'points' || name === 'order' ? Number(value) : value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    try {
      await exercisesAPI.updateExercise(exerciseUUID, exercise);
      navigate(`/lessons/${exercise.lesson_uuid}`);
    } catch (err) {
      console.error('Ошибка при обновлении упражнения:', err);
      setError('Ошибка при обновлении упражнения');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="admin-panel">
      <h2>Редактировать упражнение</h2>

      {error && <p className="error-message">{error}</p>}

      <form className='admin-form' onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Название упражнения</label>
          <input
            type="text"
            name="title"
            value={exercise.title}
            onChange={handleChange}
            required
          />
        </div>

        <div className="form-group">
          <label>Описание</label>
          <textarea
            name="description"
            value={exercise.description}
            onChange={handleChange}
            rows={8}
            required
          ></textarea>
        </div>

        <div className="form-group">
          <label>Очки</label>
          <input
            type="number"
            name="points"
            value={exercise.points}
            onChange={handleChange}
            required
          />
        </div>

        <div className="form-group">
          <label>Порядок</label>
          <input
            type="number"
            name="order"
            value={exercise.order}
            onChange={handleChange}
            required
          />
        </div>

        <button className="submit-button" type="submit" disabled={loading}>
          {loading ? 'Сохраняем...' : 'Обновить упражнение'}
        </button>
      </form>
    </div>
  );
};
