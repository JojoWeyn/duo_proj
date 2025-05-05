import './App.css';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Auth from './components/pages/Auth';
import Layout from './components/pages/Layout';
import Security from './components/pages/Security';
import { CourseList } from './components/Courses/CourseList';
import { LessonList } from './components/Courses/LessonList';
import { ExercisesList } from './components/Courses/ExerciseList';
import { QuestionList } from './components/Courses/QuestionList';
import { CourseCreate } from './components/Create/CourseCreate';
import { LessonCreate } from './components/Create/LessonCreate';
import { ExerciseCreate } from './components/Create/ExerciseCreate';
import { QuestionCreate } from './components/Create/QuestionCreate';
import UserList from './components/Users/UserList';
import UserDetail from './components/Users/UserDetail';
import PrivateRoute from './components/PrivateRoute';
import { CourseUpdate } from './components/Update/CourseUpdate';
import { LessonUpdate } from './components/Update/LessonUpdate';
import { ExerciseUpdate } from './components/Update/ExerciseUpdate';
import {FileList} from './components/Files/FileList';
import { AchievementList } from './components/Achievements/AchievementList';
import { AchievementCreate } from './components/Create/AchievementCreate';
import { AchievementUpdate } from './components/Update/AchievementUpdate';
import { QuestionUpdate } from './components/Update/QuestionUpdate';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Navigate to="/courses" replace />} />
        <Route path="/login" element={<Auth />} />

        <Route element={<PrivateRoute />}>
          <Route element={<Layout />}>
            <Route path="/courses" element={<CourseList />} />
            <Route path="/courses/:uuid" element={<LessonList />} />
            <Route path="/lessons/:uuid" element={<ExercisesList />} />
            <Route path="/exercises/:uuid" element={<QuestionList />} />
            <Route path="/questions/:uuid/update" element={<QuestionUpdate />} />
            <Route path="/course/create" element={<CourseCreate />} />
            <Route path="/courses/:courseUUID/lesson/create" element={<LessonCreate />} />
            <Route path="/lessons/:lessonUUID/exercise/create" element={<ExerciseCreate />} />
            <Route path="/exercises/:exerciseUUID/question/create" element={<QuestionCreate />} />
            <Route path="/courses/:courseUUID/update" element={<CourseUpdate/>} />
            <Route path="/lessons/:lessonUUID/update" element={<LessonUpdate/>}/>
            <Route path="/exercises/:exerciseUUID/update" element={<ExerciseUpdate/>}/>
            <Route path="/users" element={<UserList />} />
            <Route path="/users/:uuid" element={<UserDetail />} />
            <Route path="/files" element={<FileList />} />
            <Route path="/security" element={<Security />} />
            <Route path="/achievements" element={<AchievementList />} />
            <Route path="/achievements/create" element={<AchievementCreate />} />
            <Route path="/achievements/:id/update" element={<AchievementUpdate />} />
          </Route>
        </Route>
      </Routes>
    </Router>
  );
}

export default App;
