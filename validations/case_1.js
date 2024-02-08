import http from 'k6/http'
import { check } from 'k6';

export const options = {
  scenarios: {
    contacts: {
      executor: 'constant-vus',
      duration: '30s',
      vus: 25,
    },
  },
};

export default function () {
    const term = 'sky'
    let url = `https://group-${__ENV.GROUP_NUMBER}--${__ENV.APP_NAME}.furyapps.io/search?term=${term}`
    if (__ENV.ENV === "local")  {
      url = `http://localhost:8080/search?term=${term}`
    }
    const res = http.get(url, {
      headers: {
        'x-tiger-token': `Bearer ${__ENV.FURY_TOKEN}`,
      },
    })
    check(res, { 'status was 200': (r) => r.status == 200 });
}