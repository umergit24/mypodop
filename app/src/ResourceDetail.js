import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

function ResourceDetail() {
  const { resourceType, name } = useParams();
  const [resource, setResource] = useState(null);

  useEffect(() => {
    async function fetchResourceDetail() {
      try {
        const response = await fetch(`http://localhost:8080/details/${resourceType}/${name}`, {
          headers: {
            'Accept': 'text/yaml',
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
  }, [resourceType, name]);

  return (
    <div>
      <h1>{resourceType} Manifest</h1>
      <pre>{resource}</pre>
    </div>
  );
}

export default ResourceDetail;
