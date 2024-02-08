import http from 'k6/http'
import { check } from 'k6';

export const options = {
  scenarios: {
    contacts: {
      executor: 'constant-vus',
      duration: '30s',
      vus: 10,
    },
  },
};

export default function () {
    let url = `https://group-${__ENV.GROUP_NUMBER}--${__ENV.APP_NAME}.furyapps.io/report`
    if (__ENV.ENV === "local")  {
      url = 'http://localhost:8080/report'
    }
    const res = http.get(url, {
      headers: {
        'x-tiger-token': `Bearer ${__ENV.FURY_TOKEN}`,
      },
    })
    check(res, { 'status was 200': (r) => r.status == 200 });
}