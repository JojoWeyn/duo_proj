import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { coursesAPI } from '../../api/api';

export const CourseUpdate = () => {
  const { courseUUID } = useParams();
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const [course, setCourse] = useState({
    title: '',
    description: '',
    type_id: 1,
    difficulty_id: 1,
  });

  const fetchCourse = async () => {
    try {
      const res = await coursesAPI.getCourses();
      const found = res.data.find(c => c.uuid === courseUUID);
      if (found) setCourse(found);
    } catch (err) {
      console.error('Ошибка при загрузке курса', err);
    }
  };

  useEffect(() => {
    fetchCourse();
  }, [courseUUID]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setCourse(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const payload = {
        uuid: courseUUID,
        title: course.title,
        description: course.description,
        type_id: parseInt(course.type_id),
        difficulty_id: parseInt(course.difficulty_id),
      };
  
      await coursesAPI.updateCourse(courseUUID, payload);
      navigate("/courses"); 
    } catch (err) {
      setError("Ошибка при создании курса!");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="admin-panel">
      <h2>Редактировать курс</h2>
      <form className="admin-form" onSubmit={handleSubmit}>
        {error && <p className="error-message">{error}</p>}
        <div className="form-group">
          <label>Название курса</label>
          <input
  type="text"
  name="title"
  value={course.title}
  onChange={handleChange}
  required
/>

<textarea
  name="description"
  value={course.description}
  onChange={handleChange}
  rows={8}
  required
></textarea>

<select
  name="type_id"
  value={course.type_id}
  onChange={handleChange}
>
  <option value={1}>Закрытый</option>
  <option value={2}>Открытый</option>
</select>

<select
  name="difficulty_id"
  value={course.difficulty_id}
  onChange={handleChange}
>
  <option value={1}>Легкий</option>
  <option value={2}>Средний</option>
  <option value={3}>Сложный</option>
</select>

        </div>
        <button className="submit-button" type="submit" disabled={loading}>
          {loading ? "Обновление" : "Обновить курс"}
        </button>
      </form>
    </div>
  );
};
