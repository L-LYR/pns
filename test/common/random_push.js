import { sentence } from './lorem-ipsum.js';

function random_push(size, appId, deviceId) {
    return {
        deviceId: deviceId,
        appId: appId,
        ignoreFreqCtrl: true,
        ignoreOnlineCheck: true,
        message: {
            title: sentence(size - 5, size + 5),
            content: sentence(size - 5, size + 5),
        },
        retry: 3,
    }
}

function random_template_push(deviceId) {
    return {
        appId: 1234,
        deviceId: deviceId,
        retry: 3,
        ignoreFreqCtrl: true,
        ignoreOnlineCheck: true,
        message: {
            id: "1783891026876833792",
            params: {
                "title": {
                    pr: {
                        "n": "xxx",
                        "m": "xxx",
                    },
                },
                "content": {
                    pr: {
                        "n": "xxx",
                        "m": "xxx",
                    }
                }
            }
        }
    }
}

export { random_push, random_template_push }