import { useState, useEffect } from "react";
import { coursesAPI } from "../../api/api";
import { Link, useNavigate } from "react-router-dom";
import ImportExcelModal from "../Courses/ImportModal";
import ConfirmDeleteModal from "../Courses/ConfirmDeleteModal";

export const CourseList = () => {
  const [courses, setCourses] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const navigate = useNavigate();
  const [showModal, setShowModal] = useState(false);
  const [courseToDelete, setCourseToDelete] = useState(null);

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

  const confirmDelete = async () => {
    try {
      await coursesAPI.deleteCourse(courseToDelete);
      setCourses((prev) => prev.filter((course) => course.uuid !== courseToDelete));
      setCourseToDelete(null);
    } catch (error) {
      console.error("Ошибка при удалении курса:", error);
      alert("Ошибка при удалении курса.");
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
        <button onClick={() => setCourseToDelete(course.uuid)} className="delete-button full-width">
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
      <div className="create-button-container">
        <button 
          className="create-button"
          onClick={() => navigate("/course/create")}
        >
            + Добавить курс
        </button>

        <button
          className="create-button import"
          onClick={() => setShowModal(true)}
        >
          Импорт курса
        </button>
      </div>
      <ImportExcelModal
        show={showModal}
        handleClose={() => setShowModal(false)}
      />

      <ConfirmDeleteModal
        show={!!courseToDelete}
        onConfirm={confirmDelete}
        onCancel={() => setCourseToDelete(null)}
      />
      
      </div>
      
  
  );
};