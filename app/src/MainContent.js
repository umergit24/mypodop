import React from 'react';
import { Routes, Route } from 'react-router-dom';
import ResourceList from './ResourceList';
import ResourceDetail from './ResourceDetail';

function MainContent() {
  return (
    <div className="main-content">
      <Routes>
        <Route path="/list/:resourceType" element={<ResourceList />} />
        <Route path="/details/:resourceType/:name" element={<ResourceDetail />} />
      </Routes>
    </div>
  );
}

export default MainContent;
