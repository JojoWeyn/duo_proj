import { useState, useEffect } from 'react';
import { coursesAPI } from '../../api/api';
import { Link, useNavigate } from 'react-router-dom';

export const CourseList = () => {
  const [courses, setCourses] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    const loadCourses = async () => {
      try {
        const response = await coursesAPI.getCourses();
        setCourses(response.data);
      } catch (error) {
        console.error('Error loading courses:', error);
      }
    };
    loadCourses();
  }, []);

  return (
    <div className="course-list">
      <h1>Все курсы</h1>
      <div className="courses-container">
        {courses.map(course => (
          <div key={course.uuid} className="course-card">
            <h2>{course.title}</h2>
            <p>{course.description}</p>
            <div className="course-meta">
              <span>Сложность: {course.difficulty.title}</span>
              <button 
                onClick={() => navigate(`/courses/${course.uuid}`)}
                className="view-button"
              >
                Уроки →
              </button>
            </div>
          </div>
        ))}
      </div>
      <button 
        className="create-button"
        onClick={() => navigate('/admin/courses')}
      >
        + Добавить курс
      </button>
    </div>
  );
};