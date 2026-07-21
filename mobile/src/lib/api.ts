import AsyncStorage from '@react-native-async-storage/async-storage';

const API_BASE_URL = 'http://178.105.238.27/api';
const TOKEN_KEY = 'auth_token';

export interface ApiResponse<T> {
  data?: T;
  error?: string;
  message?: string;
}

export interface AuthState {
  token: string;
  user: {
    id: string;
    displayName: string;
    emailAddress: string;
    role?: string;
  };
}

async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const token = await getToken();
  
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` }),
      ...options.headers,
    },
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || 'Request failed');
  }

  return response.json();
}

async function getToken(): Promise<string | null> {
  try {
    return await AsyncStorage.getItem(TOKEN_KEY);
  } catch (error) {
    console.error('Error getting token:', error);
    return null;
  }
}

export async function setToken(token: string): Promise<void> {
  try {
    await AsyncStorage.setItem(TOKEN_KEY, token);
  } catch (error) {
    console.error('Error setting token:', error);
  }
}

export async function clearToken(): Promise<void> {
  try {
    await AsyncStorage.removeItem(TOKEN_KEY);
  } catch (error) {
    console.error('Error clearing token:', error);
  }
}

export async function isAuthenticated(): Promise<boolean> {
  const token = await getToken();
  return !!token;
}

// Auth endpoints
export async function login(identifier: string, password: string): Promise<AuthState> {
  const response = await request<{ token: string; user: any }>('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ identifier, password }),
  });
  
  await setToken(response.token);
  return { token: response.token, user: response.user };
}

export async function signup(phone: string, fullName: string, password: string): Promise<AuthState> {
  const response = await request<{ token: string; user: any }>('/auth/signup', {
    method: 'POST',
    body: JSON.stringify({ phone, fullName, password }),
  });
  
  await setToken(response.token);
  return { token: response.token, user: response.user };
}

// Business endpoints
export async function getBusinessListings() {
  return request('/business/listings');
}

export async function getBusinessListing(id: string) {
  return request(`/business/listings/${id}`);
}

export async function createBusinessListing(data: any) {
  return request('/business/listings', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

// Jobs endpoints
export async function getJobs() {
  return request('/business/jobs');
}

export async function createJob(data: any) {
  return request('/business/jobs', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

// Tenders endpoints
export async function getTenders() {
  return request('/business/tenders');
}

export async function createTender(data: any) {
  return request('/business/tenders', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

// Class Groups endpoints
export async function getClassGroups() {
  return request('/business/class-groups');
}

export async function createClassGroup(data: any) {
  return request('/business/class-groups', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

// Offers endpoints
export async function getOffers() {
  return request('/business/offers');
}

export async function createOffer(data: any) {
  return request('/business/offers', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

// Sponsorships endpoints
export async function getSponsorships() {
  return request('/business/sponsorships');
}

export async function createSponsorship(data: any) {
  return request('/business/sponsorships', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

// Escrow endpoints
export async function createEscrowTransaction(businessId: string, data: any) {
  return request(`/business/listings/${businessId}/escrow`, {
    method: 'POST',
    body: JSON.stringify(data),
  });
}
