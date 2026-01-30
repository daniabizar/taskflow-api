// API Configuration
const API_BASE_URL = 'http://localhost:8080/api/v1';

// Utility Functions
function getToken() {
    return localStorage.getItem('taskflow_token');
}

function removeToken() {
    localStorage.removeItem('taskflow_token');
}

function showAlert(message, type = 'error') {
    const alertDiv = document.getElementById('alert');
    alertDiv.textContent = message;
    alertDiv.className = `alert ${type}`;
    alertDiv.style.display = 'block';
    
    setTimeout(() => {
        alertDiv.style.display = 'none';
    }, 5000);
}

// Check authentication
const token = getToken();
if (!token) {
    window.location.href = '../index.html';
}

// Logout Function
function logout() {
    if (confirm('Are you sure you want to logout?')) {
        removeToken();
        window.location.href = '../index.html';
    }
}

// Fetch user profile
async function fetchProfile() {
    try {
        const response = await fetch(`${API_BASE_URL}/auth/profile`, {
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });
        
        const data = await response.json();
        
        if (data.success) {
            document.getElementById('userName').textContent = data.data.name;
        } else {
            if (response.status === 401) {
                removeToken();
                window.location.href = '../index.html';
            }
        }
    } catch (error) {
        console.error('Profile fetch error:', error);
    }
}

// Fetch statistics
async function fetchStats() {
    try {
        const response = await fetch(`${API_BASE_URL}/tasks/stats`, {
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });
        
        const data = await response.json();
        
        if (data.success) {
            document.getElementById('totalTasks').textContent = data.data.total;
            document.getElementById('completedTasks').textContent = data.data.completed;
            document.getElementById('pendingTasks').textContent = data.data.pending;
            document.getElementById('highPriorityTasks').textContent = data.data.high_priority;
        }
    } catch (error) {
        console.error('Stats fetch error:', error);
    }
}

// Fetch all tasks
let allTasks = [];
let currentFilter = 'all';

async function fetchTasks() {
    try {
        const response = await fetch(`${API_BASE_URL}/tasks`, {
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });
        
        const data = await response.json();
        
        if (data.success) {
            allTasks = data.data || [];
            displayTasks(allTasks);
        } else {
            document.getElementById('tasksList').innerHTML = 
                '<p class="empty-state">Failed to load tasks.</p>';
        }
    } catch (error) {
        console.error('Tasks fetch error:', error);
        document.getElementById('tasksList').innerHTML = 
            '<p class="empty-state">Network error. Please check if the API is running.</p>';
    }
}

// Display tasks
function displayTasks(tasks) {
    const tasksList = document.getElementById('tasksList');
    
    if (tasks.length === 0) {
        tasksList.innerHTML = '<p class="empty-state">No tasks found. Create your first task!</p>';
        return;
    }
    
    tasksList.innerHTML = tasks.map(task => `
        <div class="task-card ${task.is_completed ? 'completed' : ''}">
            <div class="task-header">
                <h3 class="task-title">${escapeHtml(task.title)}</h3>
                <div class="task-badges">
                    <span class="badge priority-${task.priority}">${task.priority.toUpperCase()}</span>
                    <span class="badge category">${task.category}</span>
                </div>
            </div>
            
            ${task.description ? `<p class="task-description">${escapeHtml(task.description)}</p>` : ''}
            
            <div class="task-meta">
                <span class="task-date">
                    ${task.due_date ? `Due: ${formatDate(task.due_date)}` : 'No due date'}
                </span>
                <div class="task-actions">
                    <button class="btn-complete" onclick="toggleTask(${task.id})">
                        ${task.is_completed ? 'Undo' : 'Complete'}
                    </button>
                    <button class="btn-delete" onclick="deleteTask(${task.id})">Delete</button>
                </div>
            </div>
        </div>
    `).join('');
}

// Filter tasks
function filterTasks(filter) {
    currentFilter = filter;
    
    // Update active button
    document.querySelectorAll('.filter-btn').forEach(btn => {
        btn.classList.remove('active');
        if (btn.dataset.filter === filter) {
            btn.classList.add('active');
        }
    });
    
    // Filter tasks
    let filteredTasks = allTasks;
    
    if (filter === 'pending') {
        filteredTasks = allTasks.filter(task => !task.is_completed);
    } else if (filter === 'completed') {
        filteredTasks = allTasks.filter(task => task.is_completed);
    } else if (filter === 'high') {
        filteredTasks = allTasks.filter(task => task.priority === 'high');
    }
    
    displayTasks(filteredTasks);
}

// Add filter event listeners
document.querySelectorAll('.filter-btn').forEach(btn => {
    btn.addEventListener('click', () => {
        filterTasks(btn.dataset.filter);
    });
});

// Create task
const createTaskForm = document.getElementById('createTaskForm');
createTaskForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const title = document.getElementById('title').value;
    const description = document.getElementById('description').value;
    const priority = document.getElementById('priority').value;
    const category = document.getElementById('category').value;
    const dueDate = document.getElementById('dueDate').value;
    
    const submitBtn = createTaskForm.querySelector('button[type="submit"]');
    submitBtn.disabled = true;
    submitBtn.textContent = 'Creating...';
    
    const taskData = {
        title,
        description,
        priority,
        category,
    };
    
    if (dueDate) {
        taskData.due_date = new Date(dueDate).toISOString();
    }
    
    try {
        const response = await fetch(`${API_BASE_URL}/tasks`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify(taskData),
        });
        
        const data = await response.json();
        
        if (data.success) {
            showAlert('Task created successfully!', 'success');
            createTaskForm.reset();
            
            // Refresh tasks and stats
            await fetchTasks();
            await fetchStats();
        } else {
            showAlert(data.error || 'Failed to create task.');
        }
    } catch (error) {
        console.error('Create task error:', error);
        showAlert('Network error. Please try again.');
    } finally {
        submitBtn.disabled = false;
        submitBtn.textContent = 'Create Task';
    }
});

// Toggle task completion
async function toggleTask(taskId) {
    try {
        const response = await fetch(`${API_BASE_URL}/tasks/${taskId}/complete`, {
            method: 'PATCH',
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });
        
        const data = await response.json();
        
        if (data.success) {
            showAlert('Task updated!', 'success');
            await fetchTasks();
            await fetchStats();
        } else {
            showAlert(data.error || 'Failed to update task.');
        }
    } catch (error) {
        console.error('Toggle task error:', error);
        showAlert('Network error. Please try again.');
    }
}

// Delete task
async function deleteTask(taskId) {
    if (!confirm('Are you sure you want to delete this task?')) {
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE_URL}/tasks/${taskId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });
        
        const data = await response.json();
        
        if (data.success) {
            showAlert('Task deleted successfully!', 'success');
            await fetchTasks();
            await fetchStats();
        } else {
            showAlert(data.error || 'Failed to delete task.');
        }
    } catch (error) {
        console.error('Delete task error:', error);
        showAlert('Network error. Please try again.');
    }
}

// Helper functions
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatDate(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diff = date - now;
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));
    
    if (days < 0) {
        return `Overdue by ${Math.abs(days)} day(s)`;
    } else if (days === 0) {
        return 'Due today';
    } else if (days === 1) {
        return 'Due tomorrow';
    } else {
        return date.toLocaleDateString('en-US', { 
            month: 'short', 
            day: 'numeric', 
            year: 'numeric' 
        });
    }
}

// Initialize dashboard
async function initDashboard() {
    await fetchProfile();
    await fetchStats();
    await fetchTasks();
}

// Start the app
initDashboard();