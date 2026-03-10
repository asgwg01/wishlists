// Управление аутентификацией
class Auth {
    static currentUser = null;
    
    static async init() {
        console.log('Auth.init() called');
        // Сначала проверяем токен в localStorage
        await this.checkAuth();
        this.updateNavigation();
        return this.currentUser;
    }
    
    static async checkAuth() {
        const token = localStorage.getItem('token');
        console.log('Checking auth, token present:', !!token);
        
        if (token) {
            try {
                // Пробуем получить информацию о пользователе с сервера
                const userData = await API.getCurrentUser();
                this.currentUser = {
                    id: userData.user_id,
                    email: userData.email,
                    name: userData.name
                };
                console.log('User authenticated:', this.currentUser);
                
                // Обновляем токен в куках для middleware
                document.cookie = `token=${token}; path=/; max-age=86400`;
                
                return true;
            } catch (e) {
                console.error('Invalid token', e);
                // Токен невалидный, удаляем его
                this.logout();
            }
        }
        return false;
    }
    
    static async login(email, password) {
        try {
            const data = await API.login(email, password);
            localStorage.setItem('token', data.token);
            if (data.refresh_token) {
                localStorage.setItem('refresh_token', data.refresh_token);
            }
            
            // Устанавливаем токен в куки для middleware
            document.cookie = `token=${data.token}; path=/; max-age=86400`;
            
            await this.checkAuth();
            this.updateNavigation();
            return { success: true };
        } catch (error) {
            console.error('Login error:', error);
            return { success: false, error: error.message };
        }
    }
    
    static async register(email, password, name) {
        try {
            await API.register(email, password, name);
            return { success: true };
        } catch (error) {
            console.error('Register error:', error);
            return { success: false, error: error.message };
        }
    }
    
    static async logout() {
        localStorage.removeItem('token');
        localStorage.removeItem('refresh_token');
        // Удаляем куку
        document.cookie = 'token=; path=/; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
        this.currentUser = null;
        this.updateNavigation();
        window.location.href = '/';
    }
    
    static updateNavigation() {
        console.log('updateNavigation() called');
        const headerPlaceholder = document.getElementById('header-placeholder');
        if (!headerPlaceholder) {
            console.error('Header placeholder not found');
            return;
        }
        
        const currentPath = window.location.pathname;
        
        if (this.currentUser) {
            console.log('Updating navigation for authenticated user');
            headerPlaceholder.innerHTML = `
                <header>
                    <nav>
                        <div class="logo">
                            <a href="/">🎁 Wishlist</a>
                        </div>
                        <div class="nav-links">
                            <a href="/my_wishlists" class="nav-btn ${currentPath === '/my_wishlists' ? 'active' : ''}">
                                <span>📋</span> Мои вишлисты
                            </a>
                            <a href="/browse_wishlists" class="nav-btn ${currentPath === '/browse_wishlists' ? 'active' : ''}">
                                <span>🔍</span> Найти вишлисты
                            </a>
                            <div class="user-email">
                                ${this.currentUser.name || this.currentUser.email}
                            </div>
                            <button id="logoutBtn" class="nav-btn nav-btn-logout">
                                <span>🚪</span> Выйти
                            </button>
                        </div>
                    </nav>
                </header>
            `;
            
            const logoutBtn = document.getElementById('logoutBtn');
            if (logoutBtn) {
                logoutBtn.addEventListener('click', () => this.logout());
            }
        } else {
            console.log('Updating navigation for anonymous user');
            headerPlaceholder.innerHTML = `
                <header>
                    <nav>
                        <div class="logo">
                            <a href="/">🎁 Wishlist</a>
                        </div>
                        <div class="nav-links">
                            <a href="/login" class="nav-btn nav-btn-login ${currentPath === '/login' ? 'active' : ''}">
                                <span>🔐</span> Войти
                            </a>
                            <a href="/register" class="nav-btn nav-btn-register ${currentPath === '/register' ? 'active' : ''}">
                                <span>📝</span> Регистрация
                            </a>
                        </div>
                    </nav>
                </header>
            `;
        }
    }
    
    static requireAuth() {
        if (!this.currentUser) {
            // Перенаправляем на страницу входа
            window.location.href = '/login';
            return false;
        }
        return true;
    }
}