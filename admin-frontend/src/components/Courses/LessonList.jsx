import { useState, useEffect } from 'react';
import { coursesAPI } from '../../api/api';
import { useParams, useNavigate } from 'react-router-dom';
import { Link } from 'react-router-dom';

export const LessonList = () => {
  const { uuid } = useParams();
  const navigate = useNavigate();
  const [lessons, setLessons] = useState([]);

  useEffect(() => {
    const loadLessons = async () => {
      try {
        const response = await coursesAPI.getCourseContent(uuid);
        console.log('Response data:', response.data);
        setLessons(response.data);
      } catch (error) {
        console.error('Error loading lessons:', error);
      }
    };
    loadLessons();
  }, [uuid]);

  return (
    <div className="course-detail">
      <button onClick={() => navigate(-1)} className="back-button">
        ← Назад к курсу
      </button>
      <h2>Все уроки</h2>
      <div className="lessons-list">
        {lessons.map(lesson => (
          <Link 
            to={`/lessons/${lesson.uuid}`}
            key={lesson.uuid} 
            className="lesson-card"
          >
            <div className="lesson-header">
              <h3>Урок {lesson.order}</h3>
              <span className="lesson-title">{lesson.title}</span>
            </div>
            <p className="lesson-description">{lesson.description}</p>
          </Link>
        ))}
      </div>
    </div>
  );
};