import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

function ResourceDetail() {
  const { resourceType, name } = useParams();
  const [resourceDetail, setResourceDetail] = useState('');

  useEffect(() => {
    fetch(`http://localhost:8080/details/${resourceType}/${name}`)
      .then(response => response.text())
      .then(data => setResourceDetail(data));
  }, [resourceType, name]);

  return (
    <div>
      <h1>{resourceType} - {name}</h1>
      <pre>{resourceDetail}</pre>
    </div>
  );
}

export default ResourceDetail;
