// API конфигурация
const API_BASE_URL = 'http://localhost:8088/api/v1';

// Базовый класс для API запросов
class API {
    static async getUserInfo(userId) {
        return this.request(`/auth/user/${userId}`);
    }
    static async request(endpoint, options = {}) {
    const token = localStorage.getItem('token');
    
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers
    };
    
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }
    
    try {
        const response = await fetch(`${API_BASE_URL}${endpoint}`, {
            ...options,
            headers
        });
        
        // Проверяем статус ответа
        if (response.status === 401) {
            console.log('Received 401 unauthorized response');
            localStorage.removeItem('token');
            document.cookie = 'token=; path=/; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
            
            // Если мы не на странице логина, перенаправляем
            if (!window.location.pathname.includes('/login')) {
                window.location.href = '/login';
            }
            
            throw new Error('Unauthorized');
        }
        
        // Для 204 No Content (например при DELETE)
        if (response.status === 204) {
            return { success: true };
        }
        
        const data = await response.json();
        
        if (!response.ok) {
            // Используем ErrorDTO формат
            throw new Error(data.error || 'Ошибка запроса');
        }
        
        return data;
    } catch (error) {
        console.error('API Error:', error);
        throw error;
    }
}
    
    // ===== AUTH ENDPOINTS =====
    
    /**
     * Регистрация нового пользователя
     * @param {string} email - Email пользователя
     * @param {string} password - Пароль (мин 6 символов)
     * @param {string} name - Имя пользователя
     */
    static async register(email, password, name) {
        return this.request('/auth/register', {
            method: 'POST',
            body: JSON.stringify({ email, password, name })
        });
    }
    
    /**
     * Вход в систему
     * @param {string} email - Email пользователя
     * @param {string} password - Пароль
     * @returns {Promise<AuthDTO>} - Токен и информация о пользователе
     */
    static async login(email, password) {
        return this.request('/auth/login', {
            method: 'POST',
            body: JSON.stringify({ email, password })
        });
    }
    
    /**
     * Получение информации о текущем пользователе
     * @returns {Promise<UserInfoDTO>} - Информация о пользователе
     */
    static async getCurrentUser() {
        return this.request('/auth/self');
    }
    
    // ===== WISHLIST ENDPOINTS =====
    
    /**
     * Создание нового вишлиста
     * @param {CreateWishlistRequestDTO} data - Данные вишлиста
     * @returns {Promise<WishlistDTO>}
     */
    static async createWishlist(data) {
        return this.request('/wishlists', {
            method: 'POST',
            body: JSON.stringify(data)
        });
    }
    
    /**
     * Получение публичных вишлистов
     * @param {Object} params - Параметры пагинации
     * @param {number} params.page - Номер страницы
     * @param {number} params.limit - Элементов на странице
     * @returns {Promise<WishlistListDTO>}
     */
    static async getPublicWishlists(params = {}) {
        const query = new URLSearchParams(params).toString();
        return this.request(`/wishlists/public${query ? '?' + query : ''}`);
    }
    
    /**
     * Получение вишлиста по ID
     * @param {string} id - UUID вишлиста
     * @returns {Promise<WishlistDTO>}
     */
    static async getWishlist(id) {
        return this.request(`/wishlists/${id}`);
    }
    
    /**
     * Обновление вишлиста
     * @param {string} id - UUID вишлиста
     * @param {UpdateWishlistRequestDTO} data - Данные для обновления
     * @returns {Promise<WishlistDTO>}
     */
    static async updateWishlist(id, data) {
        return this.request(`/wishlists/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data)
        });
    }
    
    /**
     * Удаление вишлиста
     * @param {string} id - UUID вишлиста
     * @returns {Promise<{success: boolean}>}
     */
    static async deleteWishlist(id) {
        return this.request(`/wishlists/${id}`, {
            method: 'DELETE'
        });
    }
    
    /**
     * Получение вишлистов пользователя
     * @param {string} userId - UUID пользователя
     * @param {Object} params - Параметры пагинации
     * @returns {Promise<WishlistListDTO>}
     */
    static async getUserWishlists(userId, params = {}) {
        const query = new URLSearchParams(params).toString();
        return this.request(`/wishlists/user/${userId}${query ? '?' + query : ''}`);
    }
    
    // ===== ITEMS ENDPOINTS =====
    
    /**
     * Добавление предмета в вишлист
     * @param {string} wishlistId - UUID вишлиста
     * @param {CreateItemRequestDTO} data - Данные предмета
     * @returns {Promise<ItemDTO>}
     */
    static async createItem(wishlistId, data) {
        return this.request(`/wishlists/${wishlistId}/items`, {
            method: 'POST',
            body: JSON.stringify(data)
        });
    }
    
    /**
     * Получение списка предметов вишлиста
     * @param {string} wishlistId - UUID вишлиста
     * @param {Object} params - Параметры пагинации
     * @returns {Promise<ItemListDTO>}
     */
    static async getItems(wishlistId, params = {}) {
        const query = new URLSearchParams(params).toString();
        return this.request(`/wishlists/${wishlistId}/items${query ? '?' + query : ''}`);
    }
    
    /**
     * Получение предмета по ID
     * @param {string} itemId - UUID предмета
     * @returns {Promise<ItemDTO>}
     */
    static async getItem(itemId) {
        return this.request(`/items/${itemId}`);
    }
    
    /**
     * Обновление предмета
     * @param {string} itemId - UUID предмета
     * @param {UpdateItemRequestDTO} data - Данные для обновления
     * @returns {Promise<ItemDTO>}
     */
    static async updateItem(itemId, data) {
        return this.request(`/items/${itemId}`, {
            method: 'PUT',
            body: JSON.stringify(data)
        });
    }
    
    /**
     * Удаление предмета
     * @param {string} itemId - UUID предмета
     * @returns {Promise<{success: boolean}>}
     */
    static async deleteItem(itemId) {
        return this.request(`/items/${itemId}`, {
            method: 'DELETE'
        });
    }
    
    /**
     * Бронирование предмета
     * @param {string} itemId - UUID предмета
     * @returns {Promise<ItemDTO>}
     */
    static async bookItem(itemId) {
        return this.request(`/items/${itemId}/book`, {
            method: 'POST'
        });
    }
    
    /**
     * Снятие брони с предмета
     * @param {string} itemId - UUID предмета
     * @returns {Promise<ItemDTO>}
     */
    static async unbookItem(itemId) {
        return this.request(`/items/${itemId}/unbook`, {
            method: 'POST'
        });
    }
    
    // ===== BOOKINGS ENDPOINTS =====
    
    /**
     * Получение бронирований текущего пользователя
     * @returns {Promise<BookingListDTO>}
     */
    static async getUserBookings() {
        return this.request('/bookings');
    }
}

// Вспомогательные функции для работы с ценами (копейки -> рубли)
API.formatPrice = function(priceInKopecks) {
    if (!priceInKopecks) return 'Цена не указана';
    return (priceInKopecks / 100).toFixed(2) + ' ₽';
};

API.parsePrice = function(priceInRubles) {
    if (!priceInRubles) return 0;
    return Math.round(parseFloat(priceInRubles) * 100);
};