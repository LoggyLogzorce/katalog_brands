// Инициализация функционала страницы
document.addEventListener('DOMContentLoaded', function() {
    // Обработчики для кнопок действий
    const actionButtons = document.querySelectorAll('.action-btn');
    actionButtons.forEach(button => {
        button.addEventListener('click', function() {
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
            } else if (this.querySelector('.fa-shopping-cart')) {
                alert('Товар добавлен в корзину!');
            }
        });
    });

    // Обработчики для кнопок профиля
    const profileButtons = document.querySelectorAll('.btn');
    profileButtons.forEach(button => {
        button.addEventListener('click', function(e) {
            if (!this.href) {
                e.preventDefault();
                alert('Функционал в разработке');
            }
        });
    });
});

// Функция для создания карточки товара
function createProductCard(item) {
    const div = document.createElement('div');
    div.className = 'product-card';
    let badgeHTML = '';
    if (item.badge) {
        badgeHTML = `<span class="badge">${item.badge}</span>`;
    }
    div.innerHTML = `
        <div style="position: relative;">
          <img src="${item.imageURL}" alt="${item.name}" class="product-image">
          ${badgeHTML}
        </div>
        <div class="product-info">
          <h3 class="product-name">${item.name}</h3>
          <div class="product-price">
            ${item.price.toLocaleString('ru-RU')} ₽
          </div>
          <div class="product-actions">
            <div class="action-btn">
              <i class="${item.isFavorite ? 'fas' : 'far'} fa-heart" style="color: #FFB6C1;"></i>
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
            document.getElementById('stat-reviews').textContent = data.reviewsCount || 5;
            document.getElementById('stat-favorites').textContent = data.favorites.length || 0;
            if (data.user_data.role === 'user') {
                document.getElementById('become-creator-btn').classList.remove('hidden');
            }

            // Подставляем избранные товары
            const favGrid = document.getElementById('favorites-grid');
            favGrid.innerHTML = '';
            if (Array.isArray(data.favorites) && data.favorites.length) {
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
        })
        .catch(err => {
            console.error(err);
            document.getElementById('user-name').textContent = 'Ошибка загрузки';
            document.getElementById('user-email').textContent = '';
            document.getElementById('favorites-grid').innerHTML = '<p style="padding:20px">Ошибка загрузки избранного.</p>';
            document.getElementById('history-grid').innerHTML = '<p style="padding:20px">Ошибка загрузки истории.</p>';
        });
});