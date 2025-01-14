import http from "k6/http";
import { check, sleep } from "k6";
import { Rate } from "k6/metrics";

const errorRate = new Rate("errors");
/**
 * In this test, we are ramping up the number of VUs until error rate is less than 1%.
 */
export const writeOptions = {
  thresholds: {
    http_req_failed: [{ threshold: "rate<0.01", abortOnFail: true }], // http errors should be less than 1%, otherwise abort the test
    http_req_duration: ["p(95)<500"], // 95% of requests should be below 500ms
  },
  scenarios: {
    breaking: {
      executor: "ramping-vus",
      stages: [
        { duration: "40s", target: 100 },
        { duration: "40s", target: 180 },
        { duration: "40s", target: 200 },
        { duration: "40s", target: 240 },
        { duration: "40s", target: 280 },
        { duration: "40s", target: 350 },
        { duration: "40s", target: 400 },
        { duration: "40s", target: 450 },
        { duration: "40s", target: 500 },
        { duration: "40s", target: 600 },
        { duration: "40s", target: 700 },
        { duration: "40s", target: 800 },
        { duration: "40s", target: 900 },
        { duration: "40s", target: 1000 },
      ],
    },
  },
};

export const options = {
  thresholds: {
    // http_req_failed: [{ threshold: "rate<0.01", abortOnFail: true }], // http errors should be less than 1%, otherwise abort the test
    http_req_duration: ["p(95)<500"], // 95% of requests should be below 500ms
  },
  scenarios: {
    breaking: {
      executor: "ramping-vus",
      stages: [
        { duration: "40s", target: 100 },
        { duration: "40s", target: 180 },
        { duration: "40s", target: 200 },
        { duration: "40s", target: 240 },
        { duration: "40s", target: 280 },
        { duration: "40s", target: 350 },
        { duration: "40s", target: 400 },
        { duration: "40s", target: 450 },
        { duration: "40s", target: 500 },
        { duration: "40s", target: 600 },
        { duration: "40s", target: 700 },
        { duration: "40s", target: 800 },
        { duration: "40s", target: 900 },
        { duration: "40s", target: 1000 },
      ],
    },
  },
};


// The function that defines VU logic.
//
// See https://grafana.com/docs/k6/latest/examples/get-started-with-k6/ to learn more
// about authoring k6 scripts.
//
function createURLs() {
  const payload = {
    URL:
      "https://example.com/route/" +
      Math.floor(Math.random() * 1000) +
      Math.floor(Math.random() * 1000) +
      Math.floor(Math.random() * 1000) +
      "abcd" +
      Math.floor(Math.random() * 1000) +
      Math.floor(Math.random() * 1000) +
      "efgh" +
      Math.floor(Math.random() * 1000) +
      Math.floor(Math.random() * 1000) +
      Math.floor(Math.random() * 1000) +
      Math.floor(Math.random() * 1000),
    userId: Math.floor(Math.random() * 1000),
  };
  const res = http.post("https://urlly.app/urls", JSON.stringify(payload));
  console.log(res.body);
  const checks = check(res, {
    "is status 200": (r) => r.status === 200,
  });
  sleep(1);
  errorRate.add(!checks);
  console.log(errorRate);
}

function readTraffic() {
  const res = http.get("https://urlly.app/eaQsiS9Y");
  const checks = check(res, {
    "is status 200": (r) => r.status === 200,
  });
  sleep(1);
  errorRate.add(!checks);
  console.log(errorRate);
}

export default function () {
  readTraffic();
}

//  execution: local
//         script: load-test.js
//         output: -
//
//      scenarios: (100.00%) 1 scenario, 500 max VUs, 2m0s max duration (incl. graceful stop):
//               * default: 500 looping VUs for 1m30s (gracefulStop: 30s)
//
//
//      data_received..................: 5.6 MB 62 kB/s
//      data_sent......................: 4.4 MB 48 kB/s
//      http_req_blocked...............: avg=63.2ms   min=0s       med=1µs      max=6.88s  p(90)=1µs      p(95)=1µs
//      http_req_connecting............: avg=5.26ms   min=0s       med=0s       max=1.25s  p(90)=0s       p(95)=0s
//      http_req_duration..............: avg=292.7ms  min=220.31ms med=236ms    max=6.63s  p(90)=279.59ms p(95)=341.85ms
//        { expected_response:true }...: avg=293.9ms  min=220.31ms med=236.04ms max=6.63s  p(90)=280.2ms  p(95)=344.63ms
//      http_req_failed................: 3.89%  1302 out of 33417
//      http_req_receiving.............: avg=2.98ms   min=4µs      med=41µs     max=1.01s  p(90)=3.45ms   p(95)=6.65ms
//      http_req_sending...............: avg=97.81µs  min=12µs     med=90µs     max=1.73ms p(90)=138µs    p(95)=174µs
//      http_req_tls_handshaking.......: avg=57.51ms  min=0s       med=0s       max=6.54s  p(90)=0s       p(95)=0s
//      http_req_waiting...............: avg=289.62ms min=220.15ms med=235.57ms max=6.57s  p(90)=273.43ms p(95)=330.49ms
//      http_reqs......................: 33417  364.963811/s
//      iteration_duration.............: avg=1.35s    min=1.22s    med=1.23s    max=8.71s  p(90)=1.28s    p(95)=1.34s
//      iterations.....................: 33417  364.963811/s
//      vus............................: 107    min=107           max=500
//      vus_max........................: 500    min=500           max=500
//
//
// running (1m31.6s), 000/500 VUs, 33417 complete and 0 interrupted iterations
// default ✓ [======================================] 500 VUs  1m30s
//

// checks.........................: 100.00% 83477 out of 83477
//      data_received..................: 11 MB   28 kB/s
//      data_sent......................: 13 MB   32 kB/s
//      errors.........................: 0.00%   0 out of 83477
//      http_req_blocked...............: avg=3.02ms   min=0s       med=1µs      max=2.07s    p(90)=1µs      p(95)=1µs
//      http_req_connecting............: avg=1.43ms   min=0s       med=0s       max=1.23s    p(90)=0s       p(95)=0s
//    ✓ http_req_duration..............: avg=239.05ms min=221.05ms med=235.02ms max=1.47s    p(90)=246.36ms p(95)=249.99ms
//        { expected_response:true }...: avg=239.05ms min=221.05ms med=235.02ms max=1.47s    p(90)=246.36ms p(95)=249.99ms
//    ✓ http_req_failed................: 0.00%   0 out of 83477
//      http_req_receiving.............: avg=1.17ms   min=6µs      med=42µs     max=951.52ms p(90)=69µs     p(95)=172µs
//      http_req_sending...............: avg=105.96µs min=14µs     med=99µs     max=1.08ms   p(90)=172µs    p(95)=219µs
//      http_req_tls_handshaking.......: avg=1.59ms   min=0s       med=0s       max=1.82s    p(90)=0s       p(95)=0s
//      http_req_waiting...............: avg=237.77ms min=220.93ms med=234.76ms max=1.47s    p(90)=246.03ms p(95)=249.56ms
//      http_reqs......................: 83477   208.032251/s
//      iteration_duration.............: avg=1.24s    min=1.22s    med=1.23s    max=3.54s    p(90)=1.24s    p(95)=1.25s
//      iterations.....................: 83477   208.032251/s
//      vus............................: 125     min=3              max=499
//      vus_max........................: 500     min=500            max=500
