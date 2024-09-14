import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';

function ResourceList() {
  const { resourceType } = useParams();
  const [resourceNames, setResourceNames] = useState([]);

  useEffect(() => {
    async function fetchResourceList() {
      try {
        const response = await fetch(`http://localhost:8080/list/${resourceType}`);
        const data = await response.json();
        setResourceNames(data);
      } catch (error) {
        console.error('Error fetching resource list:', error);
      }
    }
    fetchResourceList();
  }, [resourceType]);

  return (
    <div>
      <h1>{resourceType} List</h1>
      <ul>
        {resourceNames.map((name) => (
          <li key={name}>
            <Link to={`/details/${resourceType}/${name}`}>{name}</Link>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default ResourceList;
