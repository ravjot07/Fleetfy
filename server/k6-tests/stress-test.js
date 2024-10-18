import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '1m', target: 10 },    // Ramp up to 10 users over 1 minute
        { duration: '2m', target: 50 },    // Hold at 50 users for 2 minutes
        { duration: '2m', target: 100 },   // Ramp up to 100 users over 2 minutes
        { duration: '3m', target: 200 },   // Hold at 200 users for 3 minutes
        { duration: '2m', target: 0 },     // Ramp down to 0 users over 2 minutes
    ],
};

export default function () {
    const url = 'http://localhost:8080/login';
    const payload = JSON.stringify({
        username: 'rav',  // Valid username
        password: 'rav',  // Valid password
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    let res = http.post(url, payload, params);

    // Check if the response is 200 (OK) and contains 'Login successful'
    check(res, {
        'status was 200': (r) => r.status === 200,
        'login successful': (r) => r.body.includes('Login successful'),  // Adjust this check based on your actual success message
    });

    // Simulate a short pause between requests
    sleep(1);
}
