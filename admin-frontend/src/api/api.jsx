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

// Перехватчик для обработки ответов и автоматического обновления токенов
let isRefreshing = false;
let failedQueue = [];

const processQueue = (error, token = null) => {
  failedQueue.forEach(prom => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  
  failedQueue = [];
};

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    
    // Если ошибка 401 и запрос не на обновление токена и не было попытки повторить запрос
    if (error.response?.status === 401 && 
        !originalRequest._retry && 
        originalRequest.url !== '/v1/auth/refresh') {
      
      if (isRefreshing) {
        // Если уже идет обновление токена, добавляем запрос в очередь
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then(token => {
            originalRequest.headers['Authorization'] = `Bearer ${token}`;
            return api(originalRequest);
          })
          .catch(err => Promise.reject(err));
      }
      
      originalRequest._retry = true;
      isRefreshing = true;
      
      try {
        const refreshToken = localStorage.getItem('refresh_token');
        if (!refreshToken) {
          // Если нет refresh токена, выходим из системы
          localStorage.removeItem('access_token');
          localStorage.removeItem('refresh_token');
          window.location.href = '/login';
          return Promise.reject(error);
        }
        
        const response = await authAPI.refreshToken(refreshToken);
        const { access_token, refresh_token } = response.data;
        
        localStorage.setItem('access_token', access_token);
        localStorage.setItem('refresh_token', refresh_token);
        
        api.defaults.headers.common['Authorization'] = `Bearer ${access_token}`;
        originalRequest.headers['Authorization'] = `Bearer ${access_token}`;
        
        processQueue(null, access_token);
        isRefreshing = false;
        
        return api(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError, null);
        isRefreshing = false;
        
        // Очищаем токены и перенаправляем на страницу входа
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        window.location.href = '/login';
        
        return Promise.reject(refreshError);
      }
    }
    
    return Promise.reject(error);
  }
);

export const authAPI = {
  login: (email, password) => api.post('/v1/auth/login', { email, password }),
  refreshToken: (refreshToken) => api.post('/v1/auth/refresh', { refresh_token: refreshToken }),
  checkToken: () => api.get('/v1/auth/token/status'),
  logout: () => api.post('/v1/auth/logout'),
  getVerificationCode: (email) => api.post('/v1/auth/verification/code', null, {
    params: { email }
  }),
  resetPassword: (code, email, newPassword) => api.post('/v1/auth/password/reset', { code, email, new_password: newPassword }),
  getUserInfo: () => api.get('/v1/auth/me'),
};

export const coursesAPI = {
  getCourses: () => api.get('/v1/admin/course/list'),
  getCourseContent: (courseUuid) => api.get(`/v1/admin/course/${courseUuid}/lesson`),
  createCourse: (data) => api.post('/v1/admin/course', data),
  updateCourse: (id, data) => api.patch(`/v1/admin/course/${id}`, data),
  deleteCourse: (id) => api.delete(`/v1/admin/course/${id}`),
  importExcel: (formData) => api.post("/v1/admin/course/import-excel", formData, 
    {
      headers: { "Content-Type": "multipart/form-data" },
    }),
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
  }),
  getUser: (uuid) => api.get(`/v1/users/${uuid}`),
  getUserAchievements: (uuid) => api.get(`/v1/users/achievements/${uuid}`)
};

export const achievementsAPI = {
  getList: () => api.get('/v1/admin/achievements/list'),
  create: (data) => api.post('/v1/admin/achievements/create', data),
  update: (id, data) => api.patch(`/v1/admin/achievements/${id}`, data),
  get: (uuid) => api.get(`/v1/admin/achievements/${uuid}`),
  delete: (id) => api.delete(`/v1/admin/achievements/${id}`)
};
