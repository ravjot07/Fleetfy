import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '10s', target: 0 },    // Start with 0 users
        { duration: '10s', target: 200 },  // Spike to 200 users in 10 seconds
        { duration: '30s', target: 200 },  // Hold at 200 users for 30 seconds
        { duration: '10s', target: 0 },    // Ramp down to 0 users in 10 seconds
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
        'login successful': (r) => r.body.includes('Login successful'),  // Adjust based on actual success message
    });

    // Simulate a short pause between requests
    sleep(1);
}
