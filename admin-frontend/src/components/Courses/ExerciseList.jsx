import { useState, useEffect } from 'react';
import { exercisesAPI, filesAPI, lessonsAPI } from '../../api/api';
import { useParams, useNavigate } from 'react-router-dom';
import FileAttachModal from '../Files/FileAttachModal'
import ConfirmDeleteModal from "../Courses/ConfirmDeleteModal";

export const ExercisesList = () => {
  const { uuid } = useParams();
  const navigate = useNavigate();
  const [exercises, setExercises] = useState([]);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true);  // State for loading

  const [showAttachModal, setShowAttachModal] = useState(false);
  const [selectedQuestionUuid, setSelectedQuestionUuid] = useState(null);
  const [lessonTitle, setLessonTitle] = useState('');

  const [exerciseToDelete, setExerciseToDelete] = useState(null);


  const openAttachModal = (questionUuid) => {
    setSelectedQuestionUuid(questionUuid);
    setShowAttachModal(true);
  };

  const closeAttachModal = () => {
    setShowAttachModal(false);
    setSelectedQuestionUuid(null);
  };
  
  useEffect(() => {
    const loadExercises = async () => {
      try {
        setLoading(true);
        setError('');
        const lessonRes = await lessonsAPI.getLesson(uuid); // получаем урок
        setLessonTitle(lessonRes.data.title);
        document.title = `Упражнения урока "${lessonRes.data.title}"`;
        const lessonUUID = uuid; 
        const exercisesRes = await lessonsAPI.getLessonContent(lessonUUID);
    
        // Получаем мета-данные для всех упражнений
        const enrichedExercises = await Promise.all(
          exercisesRes.data.map(async (exercise) => {
            try {
              const meta = await exercisesAPI.getExerciseMeta(exercise.uuid);
    
              // Если мета-данные существуют, заменяем поля
              if (meta.data) {
                return {
                  ...exercise,
                  exercise_files: meta.data.exercise_files || exercise.exercise_files,
                };
              }
    
              return exercise;
            } catch (error) {
              console.error(`Error loading meta for exercise ${exercise.uuid}`, error);
              return exercise;
            }
          })
        );
    
        console.log('Exercises with updated exercise_files:', enrichedExercises);
        setExercises(enrichedExercises);
      } catch (error) {
        setError('Ошибка при загрузке упражнений!');
        console.error('Error loading exercises:', error);
      } finally {
        setLoading(false); // Set loading to false after data is fetched
      }
    };
  
    loadExercises();
  }, [uuid]);

  const handleUnpinFile = async (entity, uuid) => {
    const isConfirmed = window.confirm("Вы уверены, что хотите открепить этот файл?");
    if (!isConfirmed) return;
    try {
      await filesAPI.unpinFile(entity, uuid);
      setExercises((prevExercises) =>
        prevExercises.map((exercise) => {
          if (exercise.uuid === selectedQuestionUuid) {
            return {
              ...exercise,
              exercise_files: exercise.exercise_files.filter((file) => file.uuid !== uuid),
            };
          }
          return exercise;
        })
      );
    } catch (error) {
      console.error('Error unpinning file:', error);
      setError('Ошибка при откреплении файла!');
    }
  };

  const confirmDeleteExercise = async () => {
    try {
      await exercisesAPI.deleteExercise(exerciseToDelete);
      setExercises((prev) => prev.filter((e) => e.uuid !== exerciseToDelete));
      setExerciseToDelete(null);
    } catch (error) {
      console.error('Error deleting exercise:', error);
      setError('Ошибка при удалении упражнения!');
    }
  };

  return (
    <div className="card-list">
      <button onClick={() => navigate(-1)} className="back-button">
        ← Назад к урокам
      </button>
      <h2>Все упражнения</h2>
      {loading && <p>Загрузка упражнений...</p>}  {/* Loading message */}
      {error && <p className="error-message">{error}</p>}  {/* Error message */}
      <div className="courses-container">
        {exercises.map(exercise => (
          <div key={exercise.uuid}>
            <div onClick={() => navigate(`/exercises/${exercise.uuid}`)} className="card-item">
              <div className='lesson-header'>
                <h3>Упражнение {exercise.order}</h3>
                <span className="lesson-title">{exercise.title}</span>
                <button onClick={(event) => {
                  event.stopPropagation();
                  openAttachModal(exercise.uuid);
                }}>
                  📎 Прикрепить файл
                </button>
              </div>

              <p style={{ whiteSpace: 'pre-line' }}>{exercise.description.replace(/\\n|\n/g, '\n')}</p>

              {exercise?.exercise_files?.length > 0 && (
                <div className="exercise-files">
                  {exercise.exercise_files.map((file) => {
                    // Проверка типа файла, если это видео
                    if (file.file_url.endsWith('.mp4')) {
                      return (
                        <div key={file.uuid} className="file-preview">
                          <video width="100%" controls>
                            <source src={file.file_url} type="video/mp4" />
                            Ваш браузер не поддерживает воспроизведение видео.
                          </video>
                          <button onClick={(event) => {
                            event.stopPropagation();
                            handleUnpinFile('exercise', file.uuid)
                          }}>
                            Открепить
                          </button>
                        </div>
                      );
                    }
                    // Если это изображение
                    if (file.file_url.match(/\.(jpeg|jpg|gif|png|webp)$/)) {
                      return (
                        <div key={file.uuid} className="file-preview">
                          <img src={file.file_url} alt={file.title} style={{ width: '100%', height: 'auto' }} />
                          <button onClick={(event) => {
                            event.stopPropagation();
                            handleUnpinFile('exercise', file.uuid)
                          }}>
                            Открепить
                          </button>
                        </div>
                      );
                    }
                    // Если это PDF
                    if (file.file_url.endsWith('.pdf')) {
                      return (
                        <div key={file.uuid} className="file-preview">
                          <iframe
                            src={file.file_url}
                            width="100%"
                            height="500px"
                            title={file.title}
                          >
                            Этот браузер не поддерживает просмотр PDF.
                          </iframe>
                          <button onClick={(event) => {
                            event.stopPropagation();
                            handleUnpinFile('exercise', file.uuid)
                          }}>
                            Открепить
                          </button>
                        </div>
                      );
                    }
                    // В случае других типов файлов, можно просто отображать ссылку
                    return (
                      <div key={file.uuid} className="file-preview">
                        <a href={file.file_url} target="_blank" rel="noreferrer" className="exercise-file-link">
                          📎 {file.title}
                        </a>
                        <button onClick={(event) => {
                          event.stopPropagation();
                          handleUnpinFile('exercise', file.uuid)
                        }}>
                          Открепить
                        </button>
                      </div>
                    );
                  })}
                </div>
              )}

              <div className="card-meta">
                <span>Points: {exercise.points}</span>
              </div>
            </div>

            <div className="card-buttons">
              <button
                className="delete-button full-width"
                onClick={() => setExerciseToDelete(exercise.uuid)}
              >
                Удалить
              </button>
              <button onClick={() => navigate(`/exercises/${exercise.uuid}/update`)} className="edit-button full-width">
                Изменить
              </button>
            </div>
          </div>
        ))}
      </div>
      <div class="create-button-container">
      <button 
        className="create-button"
        onClick={() => navigate(`/lessons/${uuid}/exercise/create`)}
      >
        + Добавить упражнение
      </button>
      </div>

      {showAttachModal && (
        <FileAttachModal
          show={showAttachModal}
          handleClose={closeAttachModal}
          entity={"exercise"}
          entityUuid={selectedQuestionUuid}
        />
      )}

<ConfirmDeleteModal
  show={!!exerciseToDelete}
  onConfirm={confirmDeleteExercise}
  onCancel={() => setExerciseToDelete(null)}
/>
    </div>
  );
};