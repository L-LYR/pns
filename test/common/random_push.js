import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import { sentence } from './lorem-ipsum.js';

function random_push(size) {
    return {
        deviceId: randomIntBetween(0, 5000000),
        appId: 12345,
        title: sentence(size - 5, size + 5),
        content: sentence(size - 5, size + 5),
        retry: 3,
    }
}

export { random_push }