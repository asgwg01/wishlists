// Главный файл приложения
console.log('app.js loaded');

// Глобальные функции для HTML
window.showModal = function(modalId) {
    console.log('showModal called for:', modalId);
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.style.display = 'block';
    }
};

window.hideModal = function(modalId) {
    console.log('hideModal called for:', modalId);
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.style.display = 'none';
    }
};

// Загрузка футера
window.loadFooter = function() {
    console.log('loadFooter() called');
    const footerPlaceholder = document.getElementById('footer-placeholder');
    if (footerPlaceholder) {
        footerPlaceholder.innerHTML = `
            <footer>
                <p>&copy; 2026 Wishlist App.</p>
                <p style="margin-top: 0.5rem; font-size: 0.9rem; color: var(--text-muted);">
                    Тестовый веб клиент
                </p>
            </footer>
        `;
    } else {
        console.error('Footer placeholder not found');
    }
};

// Закрытие модальных окон по клику вне
document.addEventListener('click', function(event) {
    if (event.target.classList.contains('modal')) {
        event.target.style.display = 'none';
    }
});

// Обработка ошибок
window.addEventListener('error', function(e) {
    console.error('Global error:', e.error);
    if (typeof Components !== 'undefined') {
        Components.showAlert('Произошла ошибка. Пожалуйста, обновите страницу.', 'error');
    }
});

window.addEventListener('unhandledrejection', function(e) {
    console.error('Unhandled promise rejection:', e.reason);
    if (typeof Components !== 'undefined') {
        Components.showAlert('Произошла ошибка при выполнении запроса.', 'error');
    }
});