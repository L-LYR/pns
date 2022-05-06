import { update_target_request, push_request } from "../common/request.js"

export const options = {
  scenarios: {
    target: {
      executor: 'constant-arrival-rate',
      exec: 'update_target_request',
      duration: '6m20s',
      timeUnit: '1s',
      rate: 200,
      preAllocatedVUs: 20,
      maxVUs: 1000,
      gracefulStop: '3s',
    },
    push: {
      executor: 'ramping-arrival-rate',
      exec: 'push_request',
      timeUnit: '1s',
      preAllocatedVUs: 200,
      maxVUs: 3000,
      stages: [
        { duration: '10s', target: 700 },
        { duration: '3m', target: 700 },
        { duration: '3m', target: 700 },
        { duration: '10s', target: 0 },
      ],
      gracefulStop: '3s',
    }
  },
  summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(95)', 'p(99)', 'p(99.99)', 'count'],
};

export { update_target_request, push_request };