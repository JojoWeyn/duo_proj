import { useState, useEffect } from 'react';
import { coursesAPI, exercisesAPI, questionsAPI } from '../../api/api';
import { useParams, useNavigate } from 'react-router-dom';
import trashIcon from '../../assets/trash.svg';
import FileAttachModal from '../Files/FileAttachModal'

export const QuestionList = () => {
  const { uuid } = useParams();
  const navigate = useNavigate();
  const [questions, setQuestions] = useState([]);
  const [options, setOptions] = useState({});
  const [newOption, setNewOption] = useState({});
  const [newPair, setNewPair] = useState({ left_text: '', right_text: '' });

  const [showAttachModal, setShowAttachModal] = useState(false);
  const [selectedQuestionUuid, setSelectedQuestionUuid] = useState(null);
  
  useEffect(() => {
    const loadQuestions = async () => {
      try {
        const response = await exercisesAPI.getExerciseContent(uuid);
        const sortedQuestions = response.data.sort((a, b) => a.order - b.order);
  
        const fullQuestions = await Promise.all(
          sortedQuestions.map(async (question) => {
            try {
              const metaResponse = await questionsAPI.getQuestionMeta(question.uuid);
              return {
                ...question,
                ...metaResponse.data, // добавим мету в структуру question
              };
            } catch (error) {
              console.error(`Ошибка при загрузке мета для вопроса ${question.uuid}:`, error);
              return question; // если мета не загрузилась, вернём вопрос как есть
            }
          })
        );
  
        setQuestions(fullQuestions);

        const optionsResults = await Promise.all(
          fullQuestions.map(async (question) => {
            try {
              if (question.type_id === 1 || question.type_id === 2) {
                const optionsResponse = await questionsAPI.getQuestionOptions(question.uuid);
                return { uuid: question.uuid, options: optionsResponse.data };
              } else if (question.type_id === 3) {
                const matchingResponse = await questionsAPI.getMatchingPair(question.uuid);
                return { uuid: question.uuid, options: matchingResponse.data };
              }
            } catch (error) {
              console.error(`Ошибка при загрузке опций для вопроса ${question.uuid}:`, error);
              return { uuid: question.uuid, options: [] };
            }
          })
        );
  
        const optionsData = optionsResults.reduce((acc, result) => {
          acc[result.uuid] = result.options;
          return acc;
        }, {});
        
        setOptions(optionsData);
  
      } catch (error) {
        console.error('Ошибка при загрузке вопросов:', error);
      }
    };
    loadQuestions();

  }, [uuid]);

  const openAttachModal = (questionUuid) => {
    setSelectedQuestionUuid(questionUuid);
    setShowAttachModal(true);
  };

  const closeAttachModal = () => {
    setShowAttachModal(false);
    setSelectedQuestionUuid(null);
  };

  const handleNewOptionSubmit = async (questionUuid) => {
    try {
      const newOptionData = {
        text: newOption.text,
        is_correct: newOption.is_correct,
        questionUUID: questionUuid,
      };
      await questionsAPI.createQuestionOption(newOptionData);
      setNewOption({});
      loadOptions(questionUuid); // Re-fetch options
    } catch (error) {
      console.error('Ошибка при добавлении варианта ответа', error);
    }
  };

  const handleNewPairSubmit = async (questionUuid) => {
    try {
      const newPairData = {
        left_text: newPair.left_text,
        right_text: newPair.right_text,
        questionUUID: questionUuid,
      };
      await questionsAPI.createMatchingPair(newPairData);
      setNewPair({ left_text: '', right_text: '' });
      loadOptions(questionUuid); // Re-fetch pairs
    } catch (error) {
      console.error('Ошибка при добавлении пары соответствия', error);
    }
  };

  const handleDeleteQuestion = async (questionUuid) => {
    try {
      await questionsAPI.deleteQuestion(questionUuid);
      setQuestions((prevQuestions) => prevQuestions.filter((question) => question.uuid !== questionUuid));
      alert('Вопрос успешно удален');
    } catch (error) {
      console.error('Ошибка при удалении вопроса:', error);
      alert('Ошибка при удалении вопроса');
    }
  };

  const handleDeleteOption = async (questionUuid, optionUuid) => {
    try {
      await questionsAPI.deleteQuestionOption(optionUuid);
      setOptions((prevOptions) => ({
        ...prevOptions,
        [questionUuid]: prevOptions[questionUuid].filter((option) => option.uuid !== optionUuid),
      }));
      alert('Вариант ответа успешно удален');
    } catch (error) {
      console.error('Ошибка при удалении варианта ответа:', error);
      alert('Ошибка при удалении варианта ответа');
    }
  };

  const handleDeletePair = async (questionUuid, pairUuid) => {
    try {
      await questionsAPI.deleteMatchingPair(pairUuid);
      setOptions((prevOptions) => ({
        ...prevOptions,
        [questionUuid]: prevOptions[questionUuid].filter((pair) => pair.uuid !== pairUuid),
      }));
      alert('Пара успешно удалена');
    } catch (error) {
      console.error('Ошибка при удалении пары соответствия:', error);
      alert('Ошибка при удалении пары соответствия');
    }
  };

  const loadOptions = async (questionUuid) => {
    try {
      if (options[questionUuid]?.type_id === 1 || options[questionUuid]?.type_id === 2) {
        const optionsResponse = await questionsAPI.getQuestionOptions(questionUuid);
        setOptions((prev) => ({
          ...prev,
          [questionUuid]: optionsResponse.data,
        }));
      } else if (options[questionUuid]?.type_id === 3) {
        const matchingResponse = await questionsAPI.getMatchingPair(questionUuid);
        setOptions((prev) => ({
          ...prev,
          [questionUuid]: matchingResponse.data,
        }));
      }
    } catch (error) {
      console.error(`Error re-fetching options for question ${questionUuid}:`, error);
    }
  };

  return (
    <div className="course-detail">
      <button onClick={() => navigate(-1)} className="back-button">
        ← Назад к упражнению
      </button>
      <h2>Все Вопросы</h2>
      <div className="card-list">
        {questions.map((question) => (
          <div key={question.uuid} className="card-item" style={{cursor: "default"}}>
            <div className="lesson-header">
              <h3>Вопрос {question.order}</h3>
              <span className="lesson-title">{question.text}</span>
              <button className="delete-button" onClick={() => handleDeleteQuestion(question.uuid)}>
                Удалить вопрос
              </button>

              <button onClick={() => openAttachModal(question.uuid)}>
            📎 Прикрепить файл
              </button>

            </div>
            <p className="lesson-description">{question.type.title}</p>

            {question?.images?.length > 0 && (
        <div className="exercise-files">
          {question.images.map((file) => {
            if (file.image_url.endsWith('.mp4')) {
              return (
                <div key={file.uuid} className="file-preview">
                  <video width="100%" controls>
                    <source src={file.image_url} type="video/mp4" />
                    Ваш браузер не поддерживает воспроизведение видео.
                  </video>
                </div>
              );
            }
            // Если это изображение
            if (file.image_url.match(/\.(jpeg|jpg|gif|png|webp)$/)) {
              return (
                <div key={file.uuid} className="file-preview">
                  <img src={file.image_url} alt={file.title} style={{ width: '100%', height: 'auto' }} />
                </div>
              );
            }
            // Если это PDF
            if (file.image_url.endsWith('.pdf')) {
              return (
                <div key={file.uuid} className="file-preview">
                  <iframe
                    src={file.image_url}
                    width="100%"
                    height="500px"
                    title={file.title}
                  >
                    Этот браузер не поддерживает просмотр PDF.
                  </iframe>
                </div>
              );
            }
            return (
              <div key={file.uuid} className="file-preview">
                <a href={file.file_url} target="_blank" rel="noreferrer" className="exercise-file-link">
                  📎 {file.title}
                </a>
              </div>
            );
          })}
        </div>
      )}

            {question.type_id === 1 || question.type_id === 2 ? (
              <div className="options-list">
                {options[question.uuid]?.map((option) => (
                  <div
                    key={option.uuid}
                    className={`option-item ${option.is_correct ? 'correct' : 'incorrect'}`}
                  >
                     <div className="option-box">
                      {option.text} {option.is_correct ? '✔️' : '❌'}
                    </div>
                    <button
                      onClick={() => handleDeleteOption(question.uuid, option.uuid)}
                    >
                     <img src={trashIcon} alt="" className='icon-medium'/>
                    </button>
                  </div>
                ))}
                {newOption.questionUuid !== question.uuid && (
                  <button
                    className="button"
                    onClick={() => setNewOption({ questionUuid: question.uuid })}
                  >
                   + Добавить вариант
                  </button>
                )}

                {newOption.questionUuid === question.uuid && (
                  <div className="new-option-form">
                    <input
                      type="text"
                      placeholder="Текст варианта"
                      value={newOption.text || ''}
                      onChange={(e) => setNewOption({ ...newOption, text: e.target.value })}
                    />
                    <label>
                      Правильный:
                      <input
                        type="checkbox"
                        checked={newOption.is_correct || false}
                        onChange={(e) => setNewOption({ ...newOption, is_correct: e.target.checked })}
                      />
                    </label>
                    <button onClick={() => handleNewOptionSubmit(question.uuid)}>
                      Добавить вариант
                    </button>
                  </div>
                )}
              </div>
            ) : null}

            {question.type_id === 3 ? (
              
              <div className="matching-list">
                {options[question.uuid]?.map((pair) => (
                  <div key={pair.uuid} className="matching-pair">
                    <span className="left-text">{pair.left_text}</span>
                    <span className="right-text">{pair.right_text}</span>
                    <button
                      className="delete-button"
                      onClick={() => handleDeletePair(question.uuid, pair.uuid)}
                    >
                      <img src={trashIcon} alt="" className='icon-small'/>
                    </button>
                  </div>
                ))}
                {newPair.questionUuid !== question.uuid && (
                  <button
                    className="button"
                    onClick={() => setNewPair({ questionUuid: question.uuid })}
                  >
                    + Добавить пару соответствия
                  </button>
                )}

                {newPair.questionUuid === question.uuid && (
                  <div className="new-pair-form">
                    <input
                      type="text"
                      placeholder="Левая часть"
                      value={newPair.left_text || ''}
                      onChange={(e) => setNewPair({ ...newPair, left_text: e.target.value })}
                    />
                    <input
                      type="text"
                      placeholder="Правая часть"
                      value={newPair.right_text || ''}
                      onChange={(e) => setNewPair({ ...newPair, right_text: e.target.value })}
                    />
                    <button onClick={() => handleNewPairSubmit(question.uuid)}>
                      Добавить пару
                    </button>
                  </div>
                )}
              </div>
             
            ) : null}
          </div>
        ))}
      </div>
      <button className="create-button" onClick={() => navigate(`/exercises/${uuid}/question/create`)}>
        + Добавить вопрос
      </button>

      {showAttachModal && (
        <FileAttachModal
          show={showAttachModal}
          handleClose={closeAttachModal}
          entity={"question"}
          entityUuid={selectedQuestionUuid}
        />
      )}

    </div>
    
  );
};
