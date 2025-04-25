import { useState, useEffect } from "react";
import { coursesAPI } from "../../api/api";
import { Link, useNavigate } from "react-router-dom";

export const CourseList = () => {
  const [courses, setCourses] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    document.title = "Курсы";
    loadCourses();
  }, []);

  const loadCourses = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await coursesAPI.getCourses();
      setCourses(response.data);
    } catch (error) {
      console.error("Ошибка при загрузке курсов:", error);
      setError("Ошибка при загрузке курсов");
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (courseId) => {
    const isConfirmed = window.confirm("Вы уверены, что хотите удалить этот курс?");
    if (!isConfirmed) return;

    try {
      await coursesAPI.deleteCourse(courseId);
      setCourses((prev) => prev.filter((course) => course.uuid !== courseId));
    } catch (error) {
      console.error("Ошибка при удалении курса:", error);
    }
  };

  const renderCourse = (course) => (
    <div key={course.uuid}>
      <div onClick={() => navigate(`/courses/${course.uuid}`)} className="card-item">
        <h2>{course.title}</h2>
        <p>{course.description}</p>
        <div className="card-meta">
          <span>Сложность: {course.difficulty.title}</span>
          <span>Тип: {course.course_type.title}</span>
        </div>
      </div>
      <div className="card-buttons">
        <button onClick={() => handleDelete(course.uuid)} className="delete-button full-width">
          Удалить
        </button>
        <button onClick={() => navigate(`/courses/${course.uuid}/update`)} className="edit-button full-width">
          Изменить
        </button>
      </div>
    </div>
  );

  return (
    <div className="card-list">
      <h1>Все курсы</h1>
      {loading && <p>Загрузка...</p>}
      {error && <p>{error}</p>}
      {!loading && !error && courses.length === 0 && <p>Курсы не найдены</p>}
      <div className="courses-container">
        {courses.map(renderCourse)}
      </div>
      <button 
        className="create-button"
        onClick={() => navigate("/course/create")}
      >
        + Добавить курс
      </button>
    </div>
  );
};