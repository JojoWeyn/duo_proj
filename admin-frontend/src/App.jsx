import './App.css'
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Auth from './components/pages/Auth';
import Dashboard from './components/pages/Dashboard';
import { CourseList } from './components/Courses/CourseList';
import { LessonList } from './components/Courses/LessonList';
import { ExercisesList } from './components/Courses/ExerciseList';
import { QuestionList } from './components/Courses/QuestionList';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Auth />} />
        <Route path="/courses" element={<CourseList />} />
        <Route path="/" element={<Dashboard />} />
        <Route path="/courses/:uuid" element={<LessonList />} />,
        <Route path="/lessons/:uuid" element={<ExercisesList />}/>,
        <Route path="/exercises/:uuid" element={<QuestionList />}/>
      </Routes>
    </Router>
  );
}

export default App;