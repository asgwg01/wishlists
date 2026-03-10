// UI компоненты
class Components {
    static showAlert(message, type = 'info') {
        const alertDiv = document.createElement('div');
        alertDiv.className = `alert alert-${type}`;
        alertDiv.textContent = message;
        alertDiv.style.animation = 'fadeIn 0.3s';
        
        const main = document.querySelector('main');
        main.insertBefore(alertDiv, main.firstChild);
        
        setTimeout(() => {
            alertDiv.style.animation = 'fadeOut 0.3s';
            setTimeout(() => alertDiv.remove(), 300);
        }, 5000);
    }
    
    static showLoading(container) {
        container.innerHTML = `
            <div class="loading-container">
                <div class="spinner"></div>
            </div>
        `;
    }
    
    static showEmptyState(container, title, description, icon = '📭', action = null) {
        container.innerHTML = `
            <div class="empty-state">
                <div class="empty-state-icon">${icon}</div>
                <h3 class="empty-state-title">${title}</h3>
                <p class="empty-state-description">${description}</p>
                ${action ? `<button class="btn" onclick="${action}">Создать</button>` : ''}
            </div>
        `;
    }
    
    /**
     * Создание карточки вишлиста
     * @param {WishlistDTO} wishlist - Данные вишлиста
     * @param {boolean} isOwner - Является ли пользователь владельцем
     */
    static createWishlistCard(wishlist, isOwner = false) {
        const date = new Date(wishlist.created_at * 1000).toLocaleDateString();
        
        return `
            <div class="card" data-id="${wishlist.id}">
                <div class="card-header">
                    <h3 class="card-title">${this.escapeHtml(wishlist.title)}</h3>
                    ${wishlist.is_public ? '<span class="card-badge">🌐 Публичный</span>' : '<span class="card-badge">🔒 Приватный</span>'}
                </div>
                <div class="card-content">
                    ${this.escapeHtml(wishlist.description || 'Нет описания')}
                </div>
                <div class="card-footer">
                    <div class="card-meta">
                        <span>📅 ${date}</span>
                    </div>
                    <a href="/${isOwner ? 'my_wishlist' : 'wishlist_view'}/${wishlist.id}" 
                       class="btn btn-sm">Подробнее</a>
                </div>
            </div>
        `;
    }
    
    /**
     * Создание красивой карточки предмета
     * @param {ItemDTO} item - Данные предмета
     * @param {boolean} isOwner - Является ли пользователь владельцем
     * @param {string} currentUserId - ID текущего пользователя
     * @param {Object} bookedByUser - Информация о пользователе, который забронировал (опционально)
     */
    static createItemCard(item, isOwner = false, currentUserId = null, bookedByUser = null) {
    const price = item.price ? API.formatPrice(item.price) : 'Цена не указана';
    const hasImage = item.image_url && item.image_url.trim() !== '';
    
    // Определяем статус предмета
    let status = 'available';
    let statusText = 'Доступно';
    let statusIcon = '✅';
    let bookedByInfo = '';
    
    if (item.booked_by) {
        if (item.booked_by === currentUserId) {
            status = 'booked-by-me';
            statusText = 'Вы забронировали';
            statusIcon = '📌';
            bookedByInfo = '<span class="booked-by-me-label">Это вы</span>';
        } else {
            status = 'booked';
            statusText = 'Забронировано';
            statusIcon = '🔒';
            // Добавляем информацию о том, кто забронировал, если она доступна
            if (bookedByUser) {
                const bookerName = bookedByUser.name || 'Пользователь';
                const bookerEmail = bookedByUser.email ? bookedByUser.email : '';
                
                bookedByInfo = `
                    <div class="booker-info">
                        <span class="booker-avatar">${bookerName.charAt(0).toUpperCase()}</span>
                        <div class="booker-details">
                            <span class="booker-name">${this.escapeHtml(bookerName)}</span>
                            ${bookerEmail ? `<span class="booker-email">${this.escapeHtml(bookerEmail)}</span>` : ''}
                        </div>
                    </div>
                `;
            } else {
                bookedByInfo = '<div class="booker-info"><span class="booker-name">Кем-то забронировано</span></div>';
            }
        }
    }
    
    return `
        <div class="item-card" data-id="${item.id}" data-booked-by="${item.booked_by || ''}">
            ${hasImage ? `
                <div class="item-image">
                    <img src="${this.escapeHtml(item.image_url)}" 
                         alt="${this.escapeHtml(item.name)}" 
                         loading="lazy">
                </div>
            ` : `
                <div class="item-image-placeholder">
                    <span>📦</span>
                </div>
            `}
            
            <div class="item-content">
                <div class="item-header">
                    <h3 class="item-title">${this.escapeHtml(item.name)}</h3>
                    <div class="item-price-tag">
                        <span class="price-icon">💰</span>
                        <span class="price-value">${price}</span>
                    </div>
                </div>
                
                ${item.description ? `
                    <div class="item-description">
                        ${this.escapeHtml(item.description)}
                    </div>
                ` : ''}
                
                <div class="item-meta">
                    <div class="item-status-badge status-${status}">
                        <span class="status-icon">${statusIcon}</span>
                        <span class="status-text">${statusText}</span>
                        ${bookedByInfo}
                    </div>
                    
                    ${item.product_url ? `
                        <a href="${this.escapeHtml(item.product_url)}" target="_blank" class="item-link" title="Перейти к товару">
                            <span>🔗</span>
                        </a>
                    ` : ''}
                </div>
                
                <div class="item-actions">
                    ${this.getItemActions(item, isOwner, status)}
                </div>
            </div>
        </div>
    `;
}
    
    static getItemActions(item, isOwner, status) {
        if (isOwner) {
            // Для владельца показываем кнопку снятия брони, если предмет забронирован
            if (item.booked_by) {
                return `
                    <button class="btn-icon btn-warning" onclick="unbookItem('${item.id}')" title="Снять бронь">
                        <span>🔄</span>
                    </button>
                    <button class="btn-icon btn-secondary" onclick="editItem('${item.id}')" title="Редактировать">
                        <span>✏️</span>
                    </button>
                    <button class="btn-icon btn-danger" onclick="deleteItem('${item.id}')" title="Удалить">
                        <span>🗑️</span>
                    </button>
                `;
            } else {
                return `
                    <button class="btn-icon btn-secondary" onclick="editItem('${item.id}')" title="Редактировать">
                        <span>✏️</span>
                    </button>
                    <button class="btn-icon btn-danger" onclick="deleteItem('${item.id}')" title="Удалить">
                        <span>🗑️</span>
                    </button>
                `;
            }
        } else {
            if (status === 'available') {
                return `
                    <button class="btn-book" onclick="bookItem('${item.id}')">
                        <span>🎁</span> Забронировать
                    </button>
                `;
            } else if (status === 'booked-by-me') {
                return `
                    <button class="btn-unbook" onclick="unbookItem('${item.id}')">
                        <span>🔄</span> Снять бронь
                    </button>
                `;
            }
        }
        return '';
    }
    
    /**
     * Создание модального окна для предмета
     */
    static createItemModal() {
        const modal = document.createElement('div');
        modal.id = 'itemModal';
        modal.className = 'modal';
        modal.innerHTML = `
            <div class="modal-content">
                <div class="modal-header">
                    <h2 id="itemModalTitle">Добавить предмет</h2>
                    <span class="close" onclick="hideModal('itemModal')">&times;</span>
                </div>
                <form id="itemForm">
                    <input type="hidden" id="itemId">
                    <div class="form-group">
                        <label for="itemName">Название *</label>
                        <input type="text" id="itemName" required 
                               maxlength="80" placeholder="Например: Книга">
                    </div>
                    <div class="form-group">
                        <label for="itemDescription">Описание</label>
                        <textarea id="itemDescription" maxlength="1000" 
                                  placeholder="Подробное описание..." rows="3"></textarea>
                    </div>
                    <div class="form-row">
                        <div class="form-group half">
                            <label for="itemPrice">Цена (₽)</label>
                            <input type="number" id="itemPrice" step="0.01" min="0" 
                                   placeholder="1500.00">
                        </div>
                        <div class="form-group half">
                            <label for="itemImageUrl">Изображение</label>
                            <input type="url" id="itemImageUrl" 
                                   placeholder="https://...">
                        </div>
                    </div>
                    <div class="form-group">
                        <label for="itemProductUrl">Ссылка на товар</label>
                        <input type="url" id="itemProductUrl" 
                               placeholder="https://example.com/product/123">
                    </div>
                    <div class="form-actions">
                        <button type="button" class="btn btn-secondary" onclick="hideModal('itemModal')">Отмена</button>
                        <button type="submit" class="btn btn-primary">💾 Сохранить</button>
                    </div>
                </form>
            </div>
        `;
        
        document.body.appendChild(modal);
    }
    
    static escapeHtml(unsafe) {
        if (!unsafe) return '';
        return unsafe
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }
}