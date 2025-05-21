import { useState, useEffect } from 'react';
import { coursesAPI, lessonsAPI } from '../../api/api';
import { useParams, useNavigate } from 'react-router-dom';
import ConfirmDeleteModal from "../Courses/ConfirmDeleteModal";

export const LessonList = () => {
  const { uuid } = useParams();
  const navigate = useNavigate();
  const [lessons, setLessons] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const [lessonToDelete, setLessonToDelete] = useState(null);

  useEffect(() => {
    const loadLessons = async () => {
      try {
        setLoading(true);
        setError(null);
        const courseRes = await coursesAPI.getCourses(); // получаем все курсы
        const course = courseRes.data.find((c) => c.uuid === uuid); // находим нужный по uuid
        if (course) { // устанавливаем заголовок
          document.title = `Уроки курса "${course.title}"`;
        }
        
        const response = await coursesAPI.getCourseContent(uuid);
        const sortedLessons = response.data.sort((a, b) => a.order - b.order);
        setLessons(sortedLessons);
      } catch (error) {
        console.error('Error loading lessons:', error);
        setError('Не удалось загрузить уроки');
      } finally {
        setLoading(false);
      }
    };
    loadLessons();
  }, [uuid]);

  const confirmDelete = async () => {
    try {
      await lessonsAPI.deleteLesson(lessonToDelete);
      setLessons((prev) => prev.filter((course) => course.uuid !== lessonToDelete));
      setLessonToDelete(null);
    } catch (error) {
      console.error("Ошибка при удалении:", error);
      alert("Ошибка при удалении.");
    }
  };

  return (
    <div className="card-list">
      <button onClick={() => navigate(`/courses`)} className="back-button">
        ← Назад к курсам
      </button>
      <h2>Все уроки</h2>
      
      {loading && <p>Загрузка уроков...</p>}
      {error && <div className="error">{error}</div>}
      {!loading && !error && lessons.length === 0 && <p>Уроки не найдены</p>}

      <div className="courses-container">
        {lessons.map(lesson => (
          <div key={lesson.uuid}>
            <div onClick={() => navigate(`/lessons/${lesson.uuid}`)} className="card-item">
              <div className="lesson-header">
                <h3>Урок {lesson.order}</h3>
                <span className="lesson-title">{lesson.title}</span>
              </div>
              <p style={{ whiteSpace: 'pre-line' }}>
                {lesson.description.replace(/\\n|\n/g, '\n')}
              </p>
            </div>
            <div className="card-buttons">
              <button 
                className="delete-button full-width" 
                onClick={() => setLessonToDelete(lesson.uuid)}
              >
                Удалить
              </button>
              <button 
                onClick={() => navigate(`/lessons/${lesson.uuid}/update`)} 
                className="edit-button full-width"
              >
                Изменить
              </button>
            </div>
          </div>
        ))}
      </div>
      <div class="create-button-container">
      <button 
        className="create-button"
        onClick={() => navigate(`/courses/${uuid}/lesson/create`)}
      >
        + Добавить урок
      </button>
      </div>
      <ConfirmDeleteModal
        show={!!lessonToDelete}
        onConfirm={confirmDelete}
        onCancel={() => setLessonToDelete(null)}
      />
    </div>
  );
};