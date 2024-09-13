import React, { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';

function ResourceList() {
  const { resourceType } = useParams();
  const [resources, setResources] = useState([]);

  useEffect(() => {
    fetch(`http://localhost:8080/list/${resourceType}`)
      .then(response => response.json())
      .then(data => setResources(data));
  }, [resourceType]);

  return (
    <div>
      <h1>{resourceType}</h1>
      <ul>
        {Object.keys(resources).map(name => (
          <li key={name}>
            <Link to={`/details/${resourceType}/${name}`}>{name}</Link>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default ResourceList;
