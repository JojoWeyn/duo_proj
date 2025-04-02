import { useState, useEffect } from 'react';
import { coursesAPI } from '../../api/api';
import { useParams, useNavigate } from 'react-router-dom';

export const ExercisesList = () => {
  const { uuid } = useParams();
  const navigate = useNavigate();
  const [exercises, setExercises] = useState([]);

  useEffect(() => {
    const loadExercises = async () => {
      try {
        const response = await coursesAPI.getLessonContent(uuid);
        setExercises(response.data);
      } catch (error) {
        console.error('Error loading exercises:', error);
        alert('Failed to load lesson content');
      }
    };
    loadExercises();
  }, [uuid]);

  return (
    <div className="course-detail">
      <h2>Все упражнения</h2>
      <div className="exercises-list">
        {exercises.map(exercise => (
          <div key={exercise.uuid} className="exercise-item">
            <h3>{exercise.title}</h3>
            <p>{exercise.description}</p>
            <div className="exercise-meta">
              <span>Points: {exercise.points}</span>
              <button
                className="view-button"
                onClick={() => navigate(`/exercises/${exercise.uuid}`)}
              >
                Вопросы →
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};