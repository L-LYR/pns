import { uuidv4, randomItem } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

const os_list = ["windows", "android", "macos", "linux"]
const brand_list = ["chrome", "huawei", "vivo", "apple", "safari", "firefox", "windows"]
const model_list = ["xxxx", "yyyy", "zzzz"]
const app_version_list = ["0.0.1", "0.2.1", "0.2.2", "0.1.0", "0.1.1"]
const push_sdk_version_list = ["0.0.1", "0.0.2", "0.0.3"]

function random_device() {
    return {
        deviceId: uuidv4(),
        os: randomItem(os_list),
        brand: randomItem(brand_list),
        model: randomItem(model_list),
        tzName: "Asia/Shanghai",
        appId: 12345,
        appVersion: randomItem(app_version_list),
        pushSDKVersion: randomItem(push_sdk_version_list),
        language: "cn",
        inAppPushStatus: 1,
        systemPushStatus: 1,
        privacyPushStatus: 1,
        businessPushStatus: null,
    }
}

export { random_device }