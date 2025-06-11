// Инициализация функционала страницы
document.addEventListener('DOMContentLoaded', function () {
    // Обработчики для кнопок действий
    const actionButtons = document.querySelectorAll('.action-btn');
    actionButtons.forEach(button => {
        button.addEventListener('click', function () {
            if (this.querySelector('.fa-heart')) {
                const icon = this.querySelector('.fa-heart');
                if (icon.classList.contains('far')) {
                    icon.classList.remove('far');
                    icon.classList.add('fas');
                    icon.style.color = '#FFB6C1';
                } else {
                    icon.classList.remove('fas');
                    icon.classList.add('far');
                    icon.style.color = '';
                }
            }
        });
    });

    // Обработчики для кнопок профиля
    const changePassword = document.getElementById('change-password-btn');
    changePassword.addEventListener('click', function (e) {
        if (!this.href) {
            e.preventDefault();
            alert('Функционал в разработке');
        }
    });

    const updateRole = document.getElementById('become-creator-btn');
    updateRole.addEventListener('click', function (e) {
        if (!this.href) {
            e.preventDefault();
            fetch('/api/v1/update_role', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    role: 'creator',
                }),
            })
                .then(res => {
                    if (!res.ok) throw new Error('Не удалось сменить роль пользователя');
                })
                .then(data => {
                    alert('Роль пользователя успешно изменена\nПовторите вход для применения изменений');
                    window.location.href = '/auth';
                })
                .catch(err => {
                    alert(err.message);
                });
        }
    });
});

// Функция для создания карточки товара
function createProductCard(item) {
    const div = document.createElement('div');
    div.className = 'product-card';
    const badgeHTML = isNewProduct(item.created_at)
        ? '<span class="badge">Новинка</span>'
        : '';

    // Собираем массив URL-ов
    const images = item.product_urls.map(obj => obj.url);
    const firstImage = images[0] || '';

    div.innerHTML = `
        <div style="position: relative;">
          <img src="/static/${firstImage}" alt="${item.name}" class="product-image"
             data-images='${JSON.stringify(images)}' data-current-index="0">
          ${badgeHTML}
        </div>
        <div class="product-info">
          <h3 class="product-name">${item.name}</h3>
          <div class="product-price">
            ${item.price.toLocaleString('ru-RU')} ₽
          </div>
          <div class="product-actions">
            <div class="action-btn favorite-btn" data-id="${item.product_id}">
                <i class="${item.is_favorite ? 'fas' : 'far'} fa-heart" style="color: #FFB6C1;"></i>
            </div>
          </div>
        </div>
      `;
    return div;
}

document.addEventListener('DOMContentLoaded', () => {
    // Один запрос ко всем данным профиля
    fetch('/api/v1/profile', {
        method: 'GET',
    })
        .then(res => {
            if (!res.ok) throw new Error('Не удалось получить данные профиля');
            return res.json();
        })
        .then(data => {
            // Заполняем информацию о пользователе
            document.getElementById('user-name').textContent = data.user_data.name;
            document.getElementById('user-email').textContent = data.user_data.email;

            if (data.view_history != null) {
                document.getElementById('stat-view-history').textContent = data.view_history.length || 0;
            } else {
                document.getElementById('stat-view-history').textContent = 0;
            }

            if (data.favorites != null) {
                document.getElementById('stat-favorites').textContent = data.favorites.length || 0;
            } else {
                document.getElementById('stat-favorites').textContent = 0;
            }

            // Подставляем избранные товары
            const favGrid = document.getElementById('favorites-grid');
            favGrid.innerHTML = '';
            if (Array.isArray(data.favorites) && data.favorites != null) {
                data.favorites.forEach(item => {
                    const card = createProductCard(item);
                    favGrid.appendChild(card);
                });
            } else {
                favGrid.innerHTML = '<p style="padding: 20px">Нет избранных товаров.</p>';
            }

            // Подставляем историю просмотров
            const histGrid = document.getElementById('history-grid');
            histGrid.innerHTML = '';
            if (Array.isArray(data.view_history) && data.view_history.length) {
                data.view_history.forEach(item => {
                    const card = createProductCard(item);
                    histGrid.appendChild(card);
                });
            } else {
                histGrid.innerHTML = '<p style="padding: 20px">Нет истории просмотров.</p>';
            }
            attachItemEventListeners();
            attachFavoriteHandlers()
        })
        .catch(err => {
            console.error(err);
            document.getElementById('user-name').textContent = 'Ошибка загрузки';
            document.getElementById('user-email').textContent = '';
            document.getElementById('favorites-grid').innerHTML = '<p style="padding:20px">Ошибка загрузки избранного.</p>';
            document.getElementById('history-grid').innerHTML = '<p style="padding:20px">Ошибка загрузки истории.</p>';
        });
});