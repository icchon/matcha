import apiClient from '../api';

export const signup = (email: string, password: string) => {
    return apiClient.post('/auth/signup', { email, password });
};

export const login = (email: string, password: string) => {
    return apiClient.post('/auth/login', { email, password });
};

export const logout = () => {
    return apiClient.post('/auth/logout');
};

export const forgotPassword = (email: string) => {
    return apiClient.post('/auth/password/forgot', { email });
};

export const resetPassword = (token: string, password: string) => {
    return apiClient.post('/auth/password/reset', { token, password });
};

// OAuth methods would also go here.
// export const googleLogin = (code: string, codeVerifier: string) => {
//     return apiClient.post('/auth/oauth/google/login', { code, code_verifier: codeVerifier });
// };

// export const githubLogin = (code: string, codeVerifier: string) => {
//     return apiClient.post('/auth/oauth/github/login', { code, code_verifier: codeVerifier });
// };
