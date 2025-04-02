import { useState, useEffect } from 'react';
import { coursesAPI } from '../../api/api';
import { useParams, useNavigate } from 'react-router-dom';
import { Link } from 'react-router-dom';

export const QuestionList = () => {
  const { uuid } = useParams();
  const navigate = useNavigate();
  const [questions, setQuestions] = useState([]);

  useEffect(() => {
    const loadQuestions = async () => {
      try {
        const response = await coursesAPI.getExerciseContent(uuid);
        setQuestions(response.data);

        // Загружаем варианты ответов параллельно
        const optionsPromises = response.data.map(async question => {
          try {
            const optionsResponse = await coursesAPI.getQuestionOptions(question.uuid);
            return { uuid: question.uuid, options: optionsResponse.data };
          } catch (error) {
            console.error(`Error loading options for question ${question.uuid}:`, error);
            return { uuid: question.uuid, options: [] };
          }
        });

        const optionsResults = await Promise.all(optionsPromises);
        const optionsData = optionsResults.reduce((acc, result) => {
          acc[result.uuid] = result.options;
          return acc;
        }, {});

        setOptions(optionsData);
      } catch (error) {
        console.error('Error loading questions:', error);
      }
    };
    loadQuestions();
  }, [uuid]);

  return (
    <div className="course-detail">
      <button onClick={() => navigate(-1)} className="back-button">
        ← Назад к упражнению
      </button>
      <h2>Все Вопросы</h2>
      <div className="lessons-list">
        {questions.map(question => (
          <Link 
            to={`/questions/${question.uuid}`}
            key={question.uuid} 
            className="lesson-card"
          >
            <div className="lesson-header">
              <h3>questions {question.order}</h3>
              <span className="lesson-title">{question.text}</span>
            </div>
            <p className="lesson-description">{question.type.title}</p>
            
            {/* Enhanced Single Choice options display */}
            {question.type_id === 1 && (
              <div className="options-list">
                {question.question_options?.length > 0 ? (
                  question.question_options.map((option, index) => (
                    <div key={index} className="option-item">
                      {option.text}
                    </div>
                  ))
                ) : (
                  <div className="no-options">No options available</div>
                )}
              </div>
            )}
            {/* Handle Matching options */}
            {question.type_id === 3 && question.matching && (
              <div className="matching-preview">
                <div className="matching-side">
                  <strong>Left side:</strong>
                  {question.matching.left_side.map((item, index) => (
                    <div key={index}>{item}</div>
                  ))}
                </div>
                <div className="matching-side">
                  <strong>Right side:</strong>
                  {question.matching.right_side.map((item, index) => (
                    <div key={index}>{item}</div>
                  ))}
                </div>
              </div>
            )}
          </Link>
        ))}
      </div>
    </div>
  );
};