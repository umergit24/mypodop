import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

function Sidebar() {
  const [resourceTypes, setResourceTypes] = useState([]);

  useEffect(() => {
    async function fetchResourceTypes() {
      try {
        const response = await fetch('http://localhost:8080/');
        const data = await response.json();
        setResourceTypes(data);
      } catch (error) {
        console.error('Error fetching resource types:', error);
      }
    }
    fetchResourceTypes();
  }, []);

  return (
    <div className="sidebar">
      <h2>Resources</h2>
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

export default Sidebar;
