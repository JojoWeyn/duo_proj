import { useState, useEffect } from "react";
import { exercisesAPI, questionsAPI } from "../../api/api";
import { useParams, useNavigate } from "react-router-dom";

export const QuestionCreate = () => {
  const { exerciseUUID } = useParams();
  const navigate = useNavigate();
  const [text, setText] = useState("");
  const [typeId, setTypeId] = useState(3);
  const [order, setOrder] = useState(1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    const fetchQuestions = async () => {
      try {
        const response = await exercisesAPI.getExerciseContent(exerciseUUID);
        const lastOrder = response.data.length ? response.data[response.data.length - 1].order : 0;
        setOrder(lastOrder + 1);
      } catch (error) {
        console.error("Error fetching questions:", error);
      }
    };
    fetchQuestions();
  }, [exerciseUUID]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    const questionData = {
      text,
      type_id: typeId,
      order,
      exercise_uuid: exerciseUUID,
    };

    try {
      await questionsAPI.createQuestion(questionData);
      navigate(-1);
    } catch (err) {
      setError("Ошибка при создании вопроса!");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="admin-panel">
      <h2>Создать новый вопрос</h2>
      <form className="admin-form" onSubmit={handleSubmit}>
        {error && <p className="error-message">{error}</p>}
        <div className="form-group">
          <label>Текст вопроса</label>
          <textarea
            value={text}
            onChange={(e) => setText(e.target.value)}
            required
          ></textarea>
        </div>
        <div className="form-group">
          <label>Тип вопроса</label>
          <select value={typeId} onChange={(e) => setTypeId(Number(e.target.value))}>
            <option value={1}>Один вариант ответа</option>
            <option value={2}>Несколько вариантов</option>
            <option value={3}>Соответствие</option>
          </select>
        </div>
        <button className="submit-button" type="submit" disabled={loading}>
          {loading ? "Создание..." : "Создать вопрос"}
        </button>
      </form>
    </div>
  );
};
