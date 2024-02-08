import http from 'k6/http'
import { check, sleep } from 'k6';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
  scenarios: {
    contacts: {
      executor: 'constant-vus',
      duration: '60s',
      vus: 1,
    },
  },
}

const productIds = [
  "23423da1-1570-4c3f-8384-2e0aa022e486",
  "1f163dcb-92fe-4503-b5bb-95c6d518e7af",
  "4c835201-99eb-4330-84f2-a8f689636f38",
  "68606754-4340-482d-94c7-fabe890a901c",
  "a75602d8-a451-4d7b-9075-c8e4929d9332",
  "02b157f5-8b2b-4391-adfe-7f918635ca05",
  "99fadf7c-279e-476b-b651-061ac0709c80",
  "3ebaed0a-7b67-4987-8fa3-ed7ebb7290e2",
]

export function setup() {
  const index = __ENV.GROUP_NUMBER * 1 - 1
  const id = productIds[index]
  console.log(`product id: ${id}`)
  return { id };
}

export default function (data) {
    let changePriceRequestURL = `https://group-${__ENV.GROUP_NUMBER}--${__ENV.APP_NAME}.furyapps.io/items/${data.id}/price-changes`
    let itemURL = `https://group-${__ENV.GROUP_NUMBER}--${__ENV.APP_NAME}.furyapps.io/items/${data.id}`
    if (__ENV.ENV === "local")  {
      changePriceRequestURL = `http://localhost:8080/items/${data.id}/price-changes`
      itemURL = `http://localhost:8080/items/${data.id}`
    }
    const payload = {
      price: randomIntBetween(10000, 20000)
    }
    const params = {
      headers: {
        'x-tiger-token': `Bearer ${__ENV.FURY_TOKEN}`,
      },
    }
    http.post(changePriceRequestURL, JSON.stringify(payload), params)

    sleep(1)
    const res = http.get(itemURL, params)
    check(res, { 'final price match': (r) => r.json().price === payload.price });
}