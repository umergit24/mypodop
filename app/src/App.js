import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import MainPage from './MainPage';
import ResourceList from './ResourceList';
import ResourceDetail from './ResourceDetail';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<MainPage />} />
        <Route path="/list/:resourceType" element={<ResourceList />} />
        <Route path="/details/:resourceType/:name" element={<ResourceDetail />} />
      </Routes>
    </Router>
  );
}

export default App;
