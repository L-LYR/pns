import http from "k6/http";
import { random_device } from "../common/random_device.js"
import { random_push } from "../common/random_push.js"
import { random_app } from "../common/app_list.js";

export const options = {
  scenarios: {
    target: {
      executor: 'constant-arrival-rate',
      exec: 'update_target_request',
      duration: '9m',
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
        { duration: '2m', target: 1000 },
        { duration: '2m', target: 1000 },
        { duration: '15s', target: 1500 },
        { duration: '15s', target: 2000 },
        { duration: '30s', target: 1000 },
        { duration: '2m', target: 1000 },
        { duration: '2m', target: 0 },
      ],
      gracefulStop: '3s',
    }
  },
  summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(95)', 'p(99)', 'p(99.99)', 'count'],
};

const inbound = "http://192.168.1.2:10086";
const bizapi = "http://192.168.1.2:10087";

export function update_target_request() {
  http.request("POST", inbound + "/target", JSON.stringify(random_device()), {
    headers: { "Content-Type": "application/json" },
  });
}

export function push_request() {
  const app = random_app()
  http.request("POST", bizapi + "/push/direct", JSON.stringify(random_push(100, app.app, app.max_device)), {
    headers: { "Content-Type": "application/json" },
  });
}

