import http from "k6/http";
import { sleep } from "k6";

export const options = {
  // A number specifying the number of VUs to run concurrently.
  vus: 500,
  // A string specifying the total duration of the test run.
  duration: "90s",
};

// The function that defines VU logic.
//
// See https://grafana.com/docs/k6/latest/examples/get-started-with-k6/ to learn more
// about authoring k6 scripts.
//
export default function () {
  const payload = {
    URL: "https://example.com/route/" + Math.floor(Math.random() * 1000) + Math.floor(Math.random() * 1000),
    userId: Math.floor(Math.random() * 1000),
  };
  http.post("https://urlly.app/urls", JSON.stringify(payload));
  sleep(1);
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
