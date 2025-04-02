import { useState, useEffect } from 'react';
import { coursesAPI, lessonsAPI } from '../../api/api';
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

  const handleDelete = async (lessonUUID) => {
    const confirmDelete = window.confirm("Вы уверены, что хотите удалить этот урок?");
    if (!confirmDelete) return;
    
    try {
      await lessonsAPI.deleteLesson(lessonUUID);
      setLessons(lessons.filter(lesson => lesson.uuid !== lessonUUID));
    } catch (error) {
      console.error('Error deleting lesson:', error);
      alert("Ошибка при удалении урока");
    }
  };

  return (
    <div className="course-detail">
      <button onClick={() => navigate(-1)} className="back-button">
        ← Назад к курсам
      </button>
      <h2>Все уроки</h2>
      <div className="lessons-list">
        {lessons.map(lesson => (
          <div key={lesson.uuid} className="lesson-card">
            <Link to={`/lessons/${lesson.uuid}`} className="lesson-content">
              <div className="lesson-header">
                <h3>Урок {lesson.order}</h3>
                <span className="lesson-title">{lesson.title}</span>
              </div>
              <p className="lesson-description" style={{ whiteSpace: 'pre-line' }}>
                {lesson.description.replace(/\n|\n/g, '\n')}
              </p>
            </Link>
            <button 
              className="delete-button" 
              onClick={() => handleDelete(lesson.uuid)}
            >
              Удалить
            </button>
          </div>
        ))}
      </div>
      <button 
        className="create-button"
        onClick={() => navigate(`/courses/${uuid}/lesson/create`)}
      >
        + Добавить урок
      </button>
    </div>
  );
};