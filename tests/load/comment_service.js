import http from 'k6/http';
import { check } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('errors');

const BASE_URL = 'http://localhost:8082';

const MIN_USER_ID = 1;
const MAX_USER_ID = 5;
const MIN_POST_ID = 1;
const MAX_POST_ID = 10;

export const options = {
  scenarios: {
    create_comments: {
      executor: 'ramping-arrival-rate',
      startRate: 0,
      timeUnit: '1s',
      preAllocatedVUs: 300,
      maxVUs: 2000,
      stages: [
        { target: 2000, duration: '30s' },
        { target: 2000, duration: '1m' },
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
  const postID = Math.floor(Math.random() * (MAX_POST_ID - MIN_POST_ID + 1)) + MIN_POST_ID;
  const userID = Math.floor(Math.random() * (MAX_USER_ID - MIN_USER_ID + 1)) + MIN_USER_ID;

  const payload = JSON.stringify({
    text: 'k6 load comment',
    user_id: userID,
    post_id: postID,
  });
  const params = { headers: { 'Content-Type': 'application/json' } };

  const res = http.post(`${BASE_URL}/api/comments`, payload, params);
  const ok = check(res, {
    'create -> 201': (r) => r.status === 201,
    'has comment_id': (r) => r.json('comment_id') !== undefined,
  });
  errorRate.add(!ok);
}