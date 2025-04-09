import { useState, useEffect } from "react";
import { coursesAPI } from "../../api/api";
import { Link, useNavigate } from "react-router-dom";

export const CourseList = () => {
  const [courses, setCourses] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    loadCourses();
  }, []);

  const loadCourses = async () => {
    try {
      const response = await coursesAPI.getCourses();
      setCourses(response.data);
    } catch (error) {
      console.error("Ошибка при загрузке курсов:", error);
    }
  };

  const handleDelete = async (courseId) => {
    const isConfirmed = window.confirm("Вы уверены, что хотите удалить этот курс?");
    if (!isConfirmed) return;

    try {
      await coursesAPI.deleteCourse(courseId);
      setCourses(courses.filter((course) => course.uuid !== courseId));
    } catch (error) {
      console.error("Ошибка при удалении курса:", error);
    }
  };

  return (
    <div className="card-list">
      <h1>Все курсы</h1>
      <div className="courses-container">
        {courses.map((course) => (
          <div>
          <div onClick={() => navigate(`/courses/${course.uuid}`)} key={course.uuid} className="card-item">
            <h2>{course.title}</h2>
            <p>{course.description}</p>
              <div className="card-meta">
                <span>Сложность: {course.difficulty.title}</span>

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
           
          
        ))}
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

