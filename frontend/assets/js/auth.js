// API Configuration
const API_BASE_URL = 'http://localhost:8080/api/v1';

// Utility Functions
function showAlert(message, type = 'error') {
    const alertDiv = document.getElementById('alert');
    alertDiv.textContent = message;
    alertDiv.className = `alert ${type}`;
    alertDiv.style.display = 'block';
    
    setTimeout(() => {
        alertDiv.style.display = 'none';
    }, 5000);
}

function setToken(token) {
    localStorage.setItem('taskflow_token', token);
}

function getToken() {
    return localStorage.getItem('taskflow_token');
}

function removeToken() {
    localStorage.removeItem('taskflow_token');
}

// Check if already logged in
if (getToken() && window.location.pathname.includes('index.html')) {
    window.location.href = 'pages/dashboard.html';
}

// Login Form Handler
const loginForm = document.getElementById('loginForm');
if (loginForm) {
    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        const submitBtn = loginForm.querySelector('button[type="submit"]');
        
        // Disable button
        submitBtn.disabled = true;
        submitBtn.textContent = 'Logging in...';
        
        try {
            const response = await fetch(`${API_BASE_URL}/auth/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ email, password }),
            });
            
            const data = await response.json();
            
            if (data.success) {
                setToken(data.data.token);
                showAlert('Login successful! Redirecting...', 'success');
                
                setTimeout(() => {
                    window.location.href = 'pages/dashboard.html';
                }, 1000);
            } else {
                showAlert(data.error || 'Login failed. Please try again.');
                submitBtn.disabled = false;
                submitBtn.textContent = 'Login';
            }
        } catch (error) {
            console.error('Login error:', error);
            showAlert('Network error. Please check if the API is running.');
            submitBtn.disabled = false;
            submitBtn.textContent = 'Login';
        }
    });
}

// Register Form Handler
const registerForm = document.getElementById('registerForm');
if (registerForm) {
    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const name = document.getElementById('name').value;
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        const submitBtn = registerForm.querySelector('button[type="submit"]');
        
        // Disable button
        submitBtn.disabled = true;
        submitBtn.textContent = 'Creating account...';
        
        try {
            const response = await fetch(`${API_BASE_URL}/auth/register`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ name, email, password }),
            });
            
            const data = await response.json();
            
            if (data.success) {
                showAlert('Account created successfully! Redirecting to login...', 'success');
                
                setTimeout(() => {
                    window.location.href = '../index.html';
                }, 2000);
            } else {
                showAlert(data.error || 'Registration failed. Please try again.');
                submitBtn.disabled = false;
                submitBtn.textContent = 'Create Account';
            }
        } catch (error) {
            console.error('Register error:', error);
            showAlert('Network error. Please check if the API is running.');
            submitBtn.disabled = false;
            submitBtn.textContent = 'Create Account';
        }
    });
}