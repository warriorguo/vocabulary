import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { SearchPage } from './pages/SearchPage';
import { WordbookPage } from './pages/WordbookPage';
import './App.css';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<SearchPage />} />
        <Route path="/wordbook" element={<WordbookPage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
