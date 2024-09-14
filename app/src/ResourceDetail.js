import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

function ResourceDetail() {
  const { resourceType, name } = useParams();
  const [resource, setResource] = useState(null);
  const [isYaml, setIsYaml] = useState(false);

  useEffect(() => {
    async function fetchResourceDetail() {
      try {
        const response = await fetch(`http://localhost:8080/details/${resourceType}/${name}`, {
          headers: {
            'Accept': isYaml ? 'text/yaml' : 'application/json',
          },
        });
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        const data = await response.text();
        setResource(data);
      } catch (error) {
        console.error('Error fetching resource detail:', error);
      }
    }

    fetchResourceDetail();
  }, [resourceType, name, isYaml]);

  const handleToggleFormat = () => setIsYaml(!isYaml);

  return (
    <div>
      <h1>{resourceType} Detail</h1>
      <button onClick={handleToggleFormat}>
        {isYaml ? 'Show as JSON' : 'Show as YAML'}
      </button>
      <pre>{resource}</pre>
    </div>
  );
}

export default ResourceDetail;
