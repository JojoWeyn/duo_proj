import { useState } from "react";
import { coursesAPI } from "../../api/api";
import { useNavigate } from "react-router-dom";

export const CourseCreate = () => {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [typeId, setTypeId] = useState(2);
  const [difficultyId, setDifficultyId] = useState(1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    const courseData = {
      title,
      description,
      type_id: typeId,
      difficulty_id: difficultyId,
    };

    try {
      await coursesAPI.createCourse(courseData);
      navigate("/courses"); 
    } catch (err) {
      setError("Ошибка при создании курса!");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="admin-panel">
      <h2>Создать новый курс</h2>
      <form className="admin-form" onSubmit={handleSubmit}>
        {error && <p className="error-message">{error}</p>}
        <div className="form-group">
          <label>Название курса</label>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
          />
        </div>
        <div className="form-group">
          <label>Описание курса</label>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            required
          ></textarea>
        </div>
        <div className="form-group">
          <label>Тип курса</label>
          <select value={typeId} onChange={(e) => setTypeId(Number(e.target.value))}>
            <option value={1}>Закрытый</option>
            <option value={2}>Открытый</option>
          </select>
        </div>
        <div className="form-group">
          <label>Сложность</label>
          <select value={difficultyId} onChange={(e) => setDifficultyId(Number(e.target.value))}>
            <option value={1}>Легкий</option>
            <option value={2}>Средний</option>
            <option value={3}>Сложный</option>
          </select>
        </div>
        <button className="submit-button" type="submit" disabled={loading}>
          {loading ? "Создание..." : "Создать курс"}
        </button>
      </form>
    </div>
  );
};
