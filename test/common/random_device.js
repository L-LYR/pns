import { uuidv4, randomItem, randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import { random_app } from './app_list.js';

const os_list = ["windows", "android", "macos", "linux"]
const brand_list = ["chrome", "huawei", "vivo", "apple", "safari", "firefox", "windows"]
const model_list = ["xxxx", "yyyy", "zzzz"]
const app_version_list = ["0.0.1", "0.2.1", "0.2.2", "0.1.0", "0.1.1"]
const push_sdk_version_list = ["0.0.1", "0.0.2", "0.0.3"]

function random_device() {
    const app = random_app()
    return {
        deviceId: randomIntBetween(0, app.max_device),
        os: randomItem(os_list),
        brand: randomItem(brand_list),
        model: randomItem(model_list),
        tzName: "Asia/Shanghai",
        appId: app.app,
        appVersion: randomItem(app_version_list),
        pushSDKVersion: randomItem(push_sdk_version_list),
        language: "cn",
        inAppPushStatus: randomIntBetween(0, 1),
        systemPushStatus: randomIntBetween(0, 1),
        privacyPushStatus: randomIntBetween(0, 1),
        businessPushStatus: null,
    }
}

export { random_device }