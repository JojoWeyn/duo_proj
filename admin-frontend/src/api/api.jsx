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
  getCourses: () => api.get('/v1/admin/course/list'),
  getCourseContent: (courseUuid) => api.get(`/v1/admin/course/${courseUuid}/lesson`),
  createCourse: (data) => api.post('/v1/admin/course', data),
  updateCourse: (id, data) => api.patch(`/v1/admin/course/${id}`, data),
  deleteCourse: (id) => api.delete(`/v1/admin/course/${id}`),
};

export const lessonsAPI = {
  getLessonContent: (lessonUuid) => api.get(`/v1/admin/lesson/${lessonUuid}/exercise`),
  getLesson: (lessonUuid) => api.get(`v1/lesson/${lessonUuid}/info`),
  createLesson: (data) => api.post('/v1/admin/lesson', data),
  updateLesson: (id, data) => api.patch(`/v1/admin/lesson/${id}`, data),
  deleteLesson: (id) => api.delete(`/v1/admin/lesson/${id}`),
};

export const exercisesAPI = {
  getExerciseContent: (exerciseUuid) => api.get(`/v1/admin/exercise/${exerciseUuid}/question`),
  createExercise: (data) => api.post('/v1/admin/exercise', data),
  updateExercise: (id, data) => api.patch(`/v1/admin/exercise/${id}`, data),
  deleteExercise: (id) => api.delete(`/v1/admin/exercise/${id}`),
  getExercise: (exerciseUuid) => api.get(`/v1/exercise/${exerciseUuid}/info`),
  getExerciseMeta: (uuid) => api.get(`/v1/exercise/${uuid}/info`)
};

export const questionsAPI = {
  getQuestionOptions: (questionUuid) => api.get(`/v1/admin/question/${questionUuid}/question-option`),
  getQuestionMeta: (uuid) => api.get(`/v1/question/${uuid}/info`),
  getMatchingPair: (questionUuid) => api.get(`/v1/admin/question/${questionUuid}/matching-pair`),
  createQuestion: (data) => api.post('/v1/admin/question', data),
  updateQuestion: (id, data) => api.patch(`/v1/admin/question/${id}`, data),
  deleteQuestion: (id) => api.delete(`/v1/admin/question/${id}`),
  createQuestionOption: (data) => api.post('/v1/admin/question-option', data),
  createMatchingPair: (data) => api.post('/v1/admin/matching-pair', data),
  deleteQuestionOption: (id) => api.delete(`/v1/admin/question-option/${id}`),
  deleteMatchingPair: (id) => api.delete(`/v1/admin/matching-pair/${id}`)
};

export const filesAPI = {
  getList: () => api.get('/v1/admin/file/list'),
  addFile: (entity, uuid, data) => api.post(`/v1/admin/file/add`, data, {
    params: {
      entity,
      uuid
    }
  }),
  unpinFile: (entity, uuid) => api.post(`/v1/admin/file/unpin`, {},{
    params: {
      entity: entity,
      uuid: uuid
    }
  }),
  uploadFile: (data) => api.post('/v1/admin/file/upload', data, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  }),
  deleteFile: (file) => api.delete('/v1/admin/file/delete', {
    params: {
      file_name: file
    }
  })
};

export const usersAPI = {
  getAll: (offset = 0, limit = 50) => api.get(`/v1/admin/users`, {
    params: {
      offset,
      limit
    }
  })
};
