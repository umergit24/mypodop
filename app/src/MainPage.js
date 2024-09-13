import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

function MainPage() {
  const [resources, setResources] = useState({});

  useEffect(() => {
    fetch('http://localhost:8080/')
      .then(response => response.json())
      .then(data => setResources(data));
  }, []);

  return (
    <div>
      <h1>Kubernetes Resources</h1>
      <ul>
        {Object.keys(resources).map(resourceType => (
          <li key={resourceType}>
            <Link to={`/list/${resourceType}`}>{resourceType}</Link>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default MainPage;
