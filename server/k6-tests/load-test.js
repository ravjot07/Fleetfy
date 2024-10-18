import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    vus: 50,         // 50 virtual users
    duration: '30s', // Run the test for 30 seconds
};

export default function () {
    const url = 'http://localhost:8080/login';  // The URL for your login route
    const payload = JSON.stringify({
        username: 'rav',  // Valid username
        password: 'rav'   // Valid password
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    // Send POST request to /login
    let res = http.post(url, payload, params);

    // Check if the response status is 200 (OK) and contains 'Login successful' (adjust this according to your actual success message)
    check(res, {
        'status was 200': (r) => r.status === 200,
        'login successful': (r) => r.body.includes('Login successful'),  // Adjust based on the actual response
    });

    // Wait for 1 second before the next request
    sleep(1);
}
