import http from "k6/http";
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';


export const options = {
  stages: [
    { duration: "1m", target: 100 },
    { duration: "1m", target: 100 },
    { duration: "1m", target: 0 },
  ],
  thresholds: {
    http_req_duration: ["p(99)<1500"],
  },
};

const url = "http://localhost:10086/target";

export default function () {
  const data = {
    deviceId: uuidv4(),
    os: "windows",
    brand: "chrome",
    model: "chrome",
    tzName: "Asia/Shanghai",
    appId: 12345,
    appVersion: "3.3.3",
    pushSDKVersion: "3.3.3",
    language: "cn",
    inAppPushStatus: 1,
    systemPushStatus: 1,
    privacyPushStatus: 1,
    businessPushStatus: null,
  };
  let res = http.request("POST", url, JSON.stringify(data), {
    headers: { "Content-Type": "application/json" },
  });
  // console.log(res.json());
}
