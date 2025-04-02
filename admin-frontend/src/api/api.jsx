import axios from 'axios';

const API_BASE_URL = 'http://37.18.102.166:3211';

const api = axios.create({
  baseURL: API_BASE_URL,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const authAPI = {
  login: (email, password) => api.post('/v1/auth/login', { email, password }),
  logout: () => api.post('/v1/auth/logout'),
};

export const coursesAPI = {
  getCourses: () => api.get('/v1/course/list'),
  getCourseContent: (courseUuid) => api.get(`/v1/course/${courseUuid}/content`),
  getLessonContent: (lessonUuid) => api.get(`/v1/lesson/${lessonUuid}/content`),
  getExerciseContent: (exerciseUuid) => api.get(`/v1/exercise/${exerciseUuid}/question`),
  getQuestionOptions: (questionUuid) => api.get(`/v1/question/${questionUuid}/info`)
};