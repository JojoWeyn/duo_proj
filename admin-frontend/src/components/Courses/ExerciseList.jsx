import { useState, useEffect } from 'react';
import { exercisesAPI, lessonsAPI } from '../../api/api';
import { useParams, useNavigate } from 'react-router-dom';

export const ExercisesList = () => {
  const { uuid } = useParams();
  const navigate = useNavigate();
  const [exercises, setExercises] = useState([]);
  const [error, setError] = useState(''); // Состояние для ошибок при удалении

  useEffect(() => {
    const loadExercises = async () => {
      try {
        const response = await lessonsAPI.getLessonContent(uuid);
        setExercises(response.data);
      } catch (error) {
        console.error('Error loading exercises:', error);
        alert('Failed to load lesson content');
      }
    };
    loadExercises();
  }, [uuid]);

  // Функция для удаления упражнения
  const handleDeleteExercise = async (exerciseUuid) => {
    if (window.confirm('Вы уверены, что хотите удалить это упражнение?')) {
      try {
        await exercisesAPI.deleteExercise(exerciseUuid);  // Используем API для удаления
        setExercises((prevExercises) => prevExercises.filter((exercise) => exercise.uuid !== exerciseUuid));  // Обновляем список упражнений
      } catch (error) {
        console.error('Error deleting exercise:', error);
        setError('Ошибка при удалении упражнения!');
      }
    }
  };

  return (
    <div className="course-detail">
      <button onClick={() => navigate(-1)} className="back-button">
        ← Назад к урокам
      </button>
      <h2>Все упражнения</h2>
      {error && <p className="error-message">{error}</p>}
      <div className="exercises-list">
        {exercises.map(exercise => (
          <div key={exercise.uuid} className="exercise-item">
            <h3>{exercise.title}</h3>
            <p style={{ whiteSpace: 'pre-line' }}>{exercise.description.replace(/\\n|\n/g, '\n')}</p>
            <div className="exercise-meta">
              <span>Points: {exercise.points}</span>
              <button
                className="view-button"
                onClick={() => navigate(`/exercises/${exercise.uuid}`)}
              >
                Вопросы →
              </button>
              {/* Кнопка удаления */}
              <button
                className="delete-button"
                onClick={() => handleDeleteExercise(exercise.uuid)}
              >
                Удалить
              </button>
            </div>
          </div>
        ))}
      </div>
      <button 
        className="create-button"
        onClick={() => navigate(`/lessons/${uuid}/exercise/create`)}
      >
        + Добавить упражнение
      </button>
    </div>
  );
};
