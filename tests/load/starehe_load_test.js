import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');

// Test configuration
export const options = {
  stages: [
    { duration: '2m', target: 100 },  // Ramp up to 100 users
    { duration: '5m', target: 500 },  // Ramp up to 500 users
    { duration: '5m', target: 1000 }, // Ramp up to 1000 users
    { duration: '10m', target: 1000 }, // Stay at 1000 users
    { duration: '5m', target: 500 },  // Ramp down to 500 users
    { duration: '2m', target: 0 },    // Ramp down to 0
  ],
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'], // 95% of requests < 500ms, 99% < 1s
    http_req_failed: ['rate<0.05'], // Error rate < 5%
    errors: ['rate<0.05'],
  },
};

const BASE_URL = __ENV.API_URL || 'http://localhost:3000';

// Helper function to generate random phone number
function randomPhone() {
  return `2547${Math.floor(Math.random() * 90000000) + 10000000}`;
}

// Helper function to generate random user data
function randomUser() {
  const firstNames = ['John', 'Jane', 'Michael', 'Sarah', 'David', 'Emily', 'James', 'Emma', 'Robert', 'Lisa'];
  const lastNames = ['Smith', 'Johnson', 'Williams', 'Brown', 'Jones', 'Garcia', 'Miller', 'Davis', 'Rodriguez', 'Martinez'];
  const houses = ['Gikubu', 'Ngala', 'Geturo', 'Shaw', 'Horsten', 'Mboya', 'Shell', 'Chaka', 'Njonjo', 'Kirkley'];
  
  return {
    phone: randomPhone(),
    full_name: `${firstNames[Math.floor(Math.random() * firstNames.length)]} ${lastNames[Math.floor(Math.random() * lastNames.length)]}`,
    class_year: Math.floor(Math.random() * 30) + 1990,
    house: houses[Math.floor(Math.random() * houses.length)],
    career: ['Engineer', 'Doctor', 'Lawyer', 'Teacher', 'Business', 'Farmer', 'Artist', 'Scientist'][Math.floor(Math.random() * 8)],
    location: ['Nairobi', 'Mombasa', 'Kisumu', 'Nakuru', 'Eldoret', 'Thika'][Math.floor(Math.random() * 6)],
  };
}

// Store auth tokens
let accessToken = '';
let refreshToken = '';
let userId = '';

export function setup() {
  // Setup: Create a test user and get auth token
  const user = randomUser();
  
  // Request OTP
  const otpRes = http.post(`${BASE_URL}/api/auth/request-otp`, JSON.stringify({
    phone: user.phone,
  }), {
    headers: { 'Content-Type': 'application/json' },
  });
  
  check(otpRes, {
    'OTP request successful': (r) => r.status === 200 || r.status === 201,
  });
  
  // Simulate OTP verification (in real test, you'd need to intercept OTP)
  sleep(1);
  
  // Signup with mock OTP
  const signupRes = http.post(`${BASE_URL}/api/auth/signup`, JSON.stringify({
    phone: user.phone,
    otp: '123456', // Mock OTP
    full_name: user.full_name,
    class_year: user.class_year,
    house: user.house,
    career: user.career,
    location: user.location,
  }), {
    headers: { 'Content-Type': 'application/json' },
  });
  
  if (signupRes.status === 200 || signupRes.status === 201) {
    const data = JSON.parse(signupRes.body);
    accessToken = data.access_token;
    refreshToken = data.refresh_token;
    userId = data.user_id;
  }
  
  return { accessToken, refreshToken, userId };
}

export default function (data) {
  const headers = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${data.accessToken}`,
  };

  // Test 1: Get own profile
  const profileRes = http.get(`${BASE_URL}/api/profiles/me`, { headers });
  check(profileRes, {
    'Profile fetch successful': (r) => r.status === 200,
    'Profile response time < 300ms': (r) => r.timings.duration < 300,
  }) || errorRate.add(1);

  sleep(1);

  // Test 2: Get feed
  const feedRes = http.get(`${BASE_URL}/api/posts/feed?limit=20`, { headers });
  check(feedRes, {
    'Feed fetch successful': (r) => r.status === 200,
    'Feed response time < 500ms': (r) => r.timings.duration < 500,
  }) || errorRate.add(1);

  sleep(2);

  // Test 3: Search profiles
  const searchRes = http.post(`${BASE_URL}/api/profiles/search`, JSON.stringify({
    search_term: '',
    limit: 20,
  }), { headers });
  check(searchRes, {
    'Search successful': (r) => r.status === 200,
    'Search response time < 500ms': (r) => r.timings.duration < 500,
  }) || errorRate.add(1);

  sleep(1);

  // Test 4: Get connections
  const connectionsRes = http.get(`${BASE_URL}/api/connections/`, { headers });
  check(connectionsRes, {
    'Connections fetch successful': (r) => r.status === 200,
  }) || errorRate.add(1);

  sleep(1);

  // Test 5: Get points balance
  const pointsRes = http.get(`${BASE_URL}/api/points/balance`, { headers });
  check(pointsRes, {
    'Points balance fetch successful': (r) => r.status === 200,
  }) || errorRate.add(1);

  sleep(1);

  // Test 6: Get notifications
  const notifRes = http.get(`${BASE_URL}/api/notifications/`, { headers });
  check(notifRes, {
    'Notifications fetch successful': (r) => r.status === 200,
  }) || errorRate.add(1);

  sleep(2);

  // Test 7: Create a post (write operation)
  const createPostRes = http.post(`${BASE_URL}/api/posts/`, JSON.stringify({
    content: `Load test post at ${new Date().toISOString()}`,
    visibility: 'connections',
  }), { headers });
  check(createPostRes, {
    'Post creation successful': (r) => r.status === 200 || r.status === 201,
    'Post creation response time < 1000ms': (r) => r.timings.duration < 1000,
  }) || errorRate.add(1);

  sleep(3);
}

export function teardown(data) {
  // Cleanup: Logout
  if (data.accessToken) {
    http.post(`${BASE_URL}/api/auth/logout`, '', {
      headers: { 
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${data.accessToken}`,
      },
    });
  }
}
