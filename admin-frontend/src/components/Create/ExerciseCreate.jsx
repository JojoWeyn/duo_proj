import { useState, useEffect } from "react";
import { coursesAPI, exercisesAPI, lessonsAPI } from "../../api/api";
import { useNavigate, useParams } from "react-router-dom";

export const ExerciseCreate = () => {
  const { lessonUUID } = useParams();
  const navigate = useNavigate();
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [points, setPoints] = useState(100);
  const [order, setOrder] = useState(1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchExercises = async () => {
      try {
        const response = await lessonsAPI.getLessonContent(lessonUUID);
        const exercises = response.data;
        setOrder(exercises.length > 0 ? exercises.length + 1 : 1);
      } catch (error) {
        console.error("Error fetching exercises:", error);
      }
    };
    fetchExercises();
  }, [lessonUUID]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    const exerciseData = {
      title,
      description,
      points,
      order,
      lessonUUID,
    };

    try {
      await exercisesAPI.createExercise(exerciseData);
      navigate(-1);
    } catch (err) {
      setError("Ошибка при создании упражнения!");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="admin-panel">
      <h2>Создать упражнение</h2>
      <form className="admin-form" onSubmit={handleSubmit}>
        {error && <p className="error-message">{error}</p>}
        <div className="form-group">
          <label>Название упражнения</label>
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
          <label>Очки</label>
          <input
            type="number"
            value={points}
            onChange={(e) => setPoints(Number(e.target.value))}
            required
          />
        </div>
        <button className="submit-button" type="submit" disabled={loading}>
          {loading ? "Создание..." : "Создать упражнение"}
        </button>
      </form>
    </div>
  );
};
