const API_URL = 'http://localhost:8080/api';

function setToken(token) {
    localStorage.setItem('token', token);
}

function getToken() {
    return localStorage.getItem('token');
}

function getUser() {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
}

function setUser(user) {
    localStorage.setItem('user', JSON.stringify(user));
}

function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = 'index.html';
}

function showTab(tab) {
    document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
    document.querySelectorAll('.form').forEach(f => f.classList.remove('active'));
    
    if (tab === 'login') {
        document.querySelector('.tab').classList.add('active');
        document.getElementById('login-form').classList.add('active');
    } else {
        document.querySelectorAll('.tab')[1].classList.add('active');
        document.getElementById('register-form').classList.add('active');
    }
}

async function register() {
    const email = document.getElementById('register-email').value;
    const password = document.getElementById('register-password').value;
    const role = document.getElementById('register-role').value;

    try {
        const response = await fetch(`${API_URL}/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password, role })
        });

        const data = await response.json();

        if (response.ok) {
            setToken(data.token);
            setUser(data.user);
            window.location.href = 'dashboard.html';
        } else {
            showError(data.error || 'Ошибка регистрации');
        }
    } catch (error) {
        showError('Ошибка соединения с сервером');
    }
}

async function login() {
    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-password').value;

    try {
        const response = await fetch(`${API_URL}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });

        const data = await response.json();

        if (response.ok) {
            setToken(data.token);
            setUser(data.user);
            window.location.href = 'dashboard.html';
        } else {
            showError(data.error || 'Неверный email или пароль');
        }
    } catch (error) {
        showError('Ошибка соединения с сервером');
    }
}

async function loadTickets() {
    const token = getToken();
    if (!token) {
        window.location.href = 'index.html';
        return;
    }

    const user = getUser();
    const roleInfo = document.getElementById('role-info');
    if (roleInfo) {
        roleInfo.innerHTML = `Вы вошли как: <strong>${user.email}</strong> (${user.role === 'student' ? 'Студент' : 'Исполнитель'})`;
    }

    try {
        const response = await fetch(`${API_URL}/tickets`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (response.status === 401) {
            logout();
            return;
        }

        const tickets = await response.json();
        displayTickets(tickets);
    } catch (error) {
        document.getElementById('tickets-list').innerHTML = '<div class="error">Ошибка загрузки</div>';
    }
}

function displayTickets(tickets) {
    const container = document.getElementById('tickets-list');
    const user = getUser();

    if (!tickets || tickets.length === 0) {
        container.innerHTML = '<div class="loading">Нет заявок</div>';
        return;
    }

    let html = '';
    tickets.forEach(ticket => {
        const statusClass = `status-${ticket.status}`;
        let statusText = '';
        switch(ticket.status) {
            case 'new': statusText = 'Новая'; break;
            case 'in_progress': statusText = 'В работе'; break;
            case 'done': statusText = 'Выполнено'; break;
            default: statusText = ticket.status;
        }

        html += `
            <div class="ticket-item">
                <div class="ticket-header">
                    <span class="ticket-title">${ticket.title}</span>
                    <span class="ticket-status ${statusClass}">${statusText}</span>
                </div>
                <div class="ticket-meta">
                    <div>📍 ${ticket.location}</div>
                    <div>📁 ${ticket.category}</div>
                    <div>🆔 #${ticket.id}</div>
                    <div>📅 ${new Date(ticket.created_at).toLocaleString()}</div>
                </div>
                <div class="ticket-description">
                    ${ticket.description}
                </div>
        `;

        if (user.role === 'executor' && ticket.status !== 'done') {
            html += `
                <div class="ticket-actions">
                    ${ticket.status === 'new' ? 
                        `<button onclick="updateStatus(${ticket.id}, 'in_progress')" class="btn-secondary">Взять в работу</button>` : 
                        `<button onclick="updateStatus(${ticket.id}, 'done')" class="btn-primary">Выполнено</button>`
                    }
                </div>
            `;
        }

        html += `</div>`;
    });

    container.innerHTML = html;
}

async function updateStatus(ticketId, status) {
    const token = getToken();
    try {
        const response = await fetch(`${API_URL}/tickets/${ticketId}/status`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ status })
        });

        if (response.ok) {
            loadTickets();
        } else {
            alert('Ошибка обновления статуса');
        }
    } catch (error) {
        alert('Ошибка соединения');
    }
}

async function createTicket() {
    console.log('📝 createTicket() вызвана');
    
    const token = getToken();
    console.log('🔑 Token:', token ? token.substring(0, 20) + '...' : 'нет токена');
    
    if (!token) {
        console.log('⛔ Нет токена, редирект на login');
        window.location.href = 'index.html';
        return;
    }

    const title = document.getElementById('ticket-title').value;
    const category = document.getElementById('ticket-category').value;
    const location = document.getElementById('ticket-location').value;
    const description = document.getElementById('ticket-description').value;

    console.log('📦 Данные формы:', { title, category, location, description });

    if (!title || !category || !location || !description) {
        console.log('⚠️ Не все поля заполнены');
        showError('Все поля обязательны для заполнения');
        return;
    }

    try {
        console.log('📡 Отправка запроса на:', `${API_URL}/tickets`);
        
        const response = await fetch(`${API_URL}/tickets`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ title, category, location, description })
        });

        console.log('📨 Статус ответа:', response.status);

        if (response.ok) {
            console.log('✅ Заявка создана, редирект на dashboard');
            window.location.href = 'dashboard.html';
        } else {
            const data = await response.json();
            console.log('❌ Ошибка от сервера:', data);
            showError(data.error || 'Ошибка создания заявки');
        }
    } catch (error) {
        console.log('💥 Критическая ошибка:', error);
        showError('Ошибка соединения с сервером');
    }
}

function showError(message) {
    const errorDiv = document.getElementById('error-message');
    if (errorDiv) {
        errorDiv.textContent = message;
        errorDiv.style.display = 'block';
        setTimeout(() => {
            errorDiv.style.display = 'none';
        }, 5000);
    }
}

document.addEventListener('DOMContentLoaded', () => {
    if (window.location.pathname.includes('dashboard') || window.location.pathname.includes('create-ticket')) {
        if (!getToken()) {
            window.location.href = 'index.html';
        }
    }
    
    if (window.location.pathname.includes('dashboard')) {
        loadTickets();
    }
});
