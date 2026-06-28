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
    ramp_to_break: {
      executor: 'ramping-arrival-rate',
      startRate: 1000,
      timeUnit: '1s',
      preAllocatedVUs: 500,
      maxVUs: 5000,
      stages: [
        { target: 2000, duration: '30s' },
        { target: 4000, duration: '30s' },
        { target: 6000, duration: '30s' },
        { target: 8000, duration: '30s' },
        { target: 10000, duration: '30s' },
      ],
    },
  },
  thresholds: {
    http_req_failed: [{ threshold: 'rate<0.02' }],
    http_req_duration: [{ threshold: 'p(95)<300', abortOnFail: true }],
  },
};

export default function () {
  const postID = Math.floor(Math.random() * (MAX_POST_ID - MIN_POST_ID + 1)) + MIN_POST_ID;
  const userID = Math.floor(Math.random() * (MAX_USER_ID - MIN_USER_ID + 1)) + MIN_USER_ID;

  const payload = JSON.stringify({ text: 'stress', user_id: userID, post_id: postID });
  const params = { headers: { 'Content-Type': 'application/json' } };

  const res = http.post(`${BASE_URL}/api/comments`, payload, params);
  errorRate.add(!check(res, { 'create -> 201': (r) => r.status === 201 }));
}