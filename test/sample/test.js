import http from "k6/http";
import { random_device } from "../common/random_device.js"
import { random_push } from "../common/random_push.js"

export const options = {
  scenarios: {
    constant_request_rate: {
      executor: 'constant-arrival-rate',
      rate: 1000,
      timeUnit: '1s',
      duration: '3m',
      preAllocatedVUs: 20,
      maxVUs: 100,
    },
  },
  summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(95)', 'p(99)', 'p(99.99)', 'count'],
};

const url = "http://192.168.1.2:10087/push/direct";

export default function () {
  http.request("POST", url, JSON.stringify(random_push(100)), {
    headers: { "Content-Type": "application/json" },
  });
}
