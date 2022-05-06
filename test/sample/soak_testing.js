import { update_target_request, push_request, template_push_request, range_push_request } from "../common/request.js"

export const options = {
    scenarios: {
        target: {
            executor: 'constant-arrival-rate',
            exec: 'update_target_request',
            duration: '1h1m',
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
            maxVUs: 1000,
            stages: [
                { duration: '30s', target: 700 },
                { duration: '1h', target: 700 },
                { duration: '30s', target: 0 },
            ],
            gracefulStop: '3s',
        },
        template: {
            executor: 'ramping-arrival-rate',
            exec: 'template_push_request',
            timeUnit: '1s',
            preAllocatedVUs: 200,
            maxVUs: 1000,
            stages: [
                { duration: '30s', target: 300 },
                { duration: '1h', target: 300 },
                { duration: '30s', target: 0 },
            ],
            gracefulStop: '3s',
        },
        range: {
            executor: 'constant-arrival-rate',
            exec: 'template_push_request',
            timeUnit: '10m',
            preAllocatedVUs: 1,
            maxVUs: 1,
            rate: 1,
            duration: '1h1m',
            gracefulStop: '3s',
        }
    },
    summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(95)', 'p(99)', 'p(99.99)', 'count'],
};


export { update_target_request, push_request, template_push_request, range_push_request };