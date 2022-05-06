import http from "k6/http";
import { random_device } from "./random_device.js"
import { random_push, random_template_push } from "./random_push.js"
import { random_app } from "./app_list.js";
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

const inbound = "http://192.168.1.2:10086";
const bizapi = "http://192.168.1.2:10087";


export function update_target_request() {
    http.request("POST", inbound + "/target", JSON.stringify(random_device()), {
        headers: { "Content-Type": "application/json" },
    });
}

export function range_push_request() {
    http.request("POST", bizapi + "/puah/template/range",
        JSON.stringify({
            "appId": 12341,
            "ignoreFreqCtrl": true,
            "message": {
                "id": 1783891026876833792,
                "params": {
                    "title": {
                        "pr": {
                            "n": "xxx",
                            "m": "xxx"
                        }
                    },
                    "content": {
                        "pr": {
                            "n": "xxx",
                            "m": "xxx"
                        }
                    }
                }
            }
        }), {
        headers: { "Content-Type": "application/json" },
    }
    )
}

export function template_push_request() {
    const deviceId = randomIntBetween(0, 1000000)
    const res = http.request("POST", bizapi + "/push/template/direct",
        JSON.stringify(random_template_push(deviceId)), {
        headers: { "Content-Type": "application/json" },
    });
    const result = res.json();
    if (result == undefined) {
        return
    }
    const payload = result.payload;
    if (payload == null) {
        return
    }
    const taskId = payload.pushTaskId;
    if (randomIntBetween(0, 9) > 6) {
        return
    }
    http.request("POST", inbound + "/log", JSON.stringify(
        {
            where: "receive",
            hint: "success",
            timestamp: Date.now(),
            appId: 1234,
            deviceId: deviceId,
            taskId: taskId,
        }
    ), {
        headers: { "Content-Type": "application/json" },
    });
    if (randomIntBetween(0, 9) > 8) {
        return
    }
    http.request("POST", inbound + "/log", JSON.stringify(
        {
            where: "show",
            hint: "success",
            timestamp: Date.now(),
            appId: 1234,
            deviceId: deviceId,
            taskId: taskId,
        }
    ), {
        headers: { "Content-Type": "application/json" },
    });
}

export function push_request() {
    const app = random_app()
    const deviceId = randomIntBetween(0, app.max_device)
    const res = http.request("POST", bizapi + "/push/direct",
        JSON.stringify(random_push(100, app.app, deviceId)), {
        headers: { "Content-Type": "application/json" },
    });
    const result = res.json();
    if (result == undefined) {
        return
    }
    const payload = result.payload;
    if (payload == null) {
        return
    }
    const taskId = payload.pushTaskId;
    if (randomIntBetween(0, 9) > 6) {
        return
    }
    http.request("POST", inbound + "/log", JSON.stringify(
        {
            where: "receive",
            hint: "success",
            timestamp: Date.now(),
            appId: app.app,
            deviceId: deviceId,
            taskId: taskId,
        }
    ), {
        headers: { "Content-Type": "application/json" },
    });
    if (randomIntBetween(0, 9) > 8) {
        return
    }
    http.request("POST", inbound + "/log", JSON.stringify(
        {
            where: "show",
            hint: "success",
            timestamp: Date.now(),
            appId: app.app,
            deviceId: deviceId,
            taskId: taskId,
        }
    ), {
        headers: { "Content-Type": "application/json" },
    });
}

