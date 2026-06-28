import http from 'k6/http';
import { check } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('errors');

const BASE_URL = 'http://localhost:8081';

export const options = {
  scenarios: {
    load: {
      executor: 'ramping-arrival-rate',
      startRate: 0,
      timeUnit: '1s',
      preAllocatedVUs: 200,
      maxVUs: 1000,
      stages: [
        { target: 1000, duration: '30s' },
        { target: 1000, duration: '1m' },
        { target: 0, duration: '15s' },
      ],
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<300'],
    errors: ['rate<0.01'],
  },
};

export default function () {
  const params = { headers: { 'Content-Type': 'application/json' } };

  // 1) создаём пост
  const payload = JSON.stringify({
    title: 'load test',
    text: 'k6 generated post',
    user_id: 1,
  });

  const createRes = http.post(`${BASE_URL}/api/posts`, payload, params);
  const createdOK = check(createRes, {
    'create -> 201': (r) => r.status === 201,
    'create has post_id': (r) => r.json('post_id') !== undefined,
  });
  errorRate.add(!createdOK);

  if (createRes.status === 201) {
    const postId = createRes.json('post_id');
    const getRes = http.get(`${BASE_URL}/api/posts/${postId}`);
    const gotOK = check(getRes, {
      'get -> 200': (r) => r.status === 200,
    });
    errorRate.add(!gotOK);
  }
}