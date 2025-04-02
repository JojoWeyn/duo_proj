import { useState, useEffect } from "react";
import { coursesAPI, lessonsAPI } from "../../api/api";
import { useNavigate, useParams } from "react-router-dom";

export const LessonCreate = () => {
  const { courseUUID } = useParams();
  const navigate = useNavigate();

  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [difficultyId, setDifficultyId] = useState(1);
  const [order, setOrder] = useState(1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchLessons = async () => {
      try {
        const response = await coursesAPI.getCourseContent(courseUUID);
        const lessons = response.data;
        if (lessons.length > 0) {
          const lastOrder = lessons[lessons.length - 1].order;
          setOrder(lastOrder + 1);
        }
      } catch (err) {
        console.error("Ошибка при загрузке уроков:", err);
      }
    };

    if (courseUUID) {
      fetchLessons();
    }
  }, [courseUUID]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    const lessonData = {
      title,
      description,
      difficulty_id: difficultyId,
      order,
      courseUUID,
    };

    console.log("Отправляем данные:", lessonData);

    try {
      await lessonsAPI.createLesson(lessonData);
      navigate(`/courses/${courseUUID}`);
    } catch (err) {
    
      setError("Ошибка при создании урока!");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="admin-panel">
      <h2>Создать новый урок</h2>
      <form className="admin-form" onSubmit={handleSubmit}>
        {error && <p className="error-message">{error}</p>}
        
        <div className="form-group">
          <label>Название урока</label>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
          />
        </div>
        
        <div className="form-group">
          <label>Описание</label>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            required
          ></textarea>
        </div>

        <div className="form-group">
          <label>Сложность</label>
          <select value={difficultyId} onChange={(e) => setDifficultyId(Number(e.target.value))}>
            <option value={1}>Легкий</option>
            <option value={2}>Средний</option>
            <option value={3}>Сложный</option>
          </select>
        </div>

        <div className="form-group">
          <label>Порядок (автоопределён)</label>
          <input type="number" value={order} disabled />
        </div>

        <button className="submit-button" type="submit" disabled={loading}>
          {loading ? "Создание..." : "Создать урок"}
        </button>
      </form>
    </div>
  );
};
