import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

function MainPage() {
  const [resourceTypes, setResourceTypes] = useState([]);

  useEffect(() => {
    async function fetchResourceTypes() {
      try {
        const response = await fetch('http://localhost:8080/');
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        const data = await response.json();
        setResourceTypes(data);
      } catch (error) {
        console.error('Error fetching resource types:', error);
      }
    }

    fetchResourceTypes();
  }, []);

  return (
    <div>
      <h1>Resource Types</h1>
      <ul>
        {resourceTypes.map((type) => (
          <li key={type}>
            <Link to={`/list/${type}`}>{type}</Link>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default MainPage;
