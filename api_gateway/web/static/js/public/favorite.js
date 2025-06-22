document.addEventListener('DOMContentLoaded', function () {
    const listContainer = document.getElementById('favorites-list');

    // Создаёт DOM-элемент для карточки товара
    function createFavoritesItem(item) {
        const div = document.createElement('div');
        div.className = 'favorites-item';

        // Собираем массив URL-ов
        const images = item.product_urls.map(obj => obj.url);
        const firstImage = images[0] || '';

        // Формируем разметку
        const badgeHTML = isNewProduct(item.created_at)
            ? '<span class="badge">Новинка</span>'
            : '';

        div.innerHTML = `
        <div style="position: relative;">
          <img src="/static/${firstImage}" alt="${item.name}" class="item-image"
             data-images='${JSON.stringify(images)}' data-current-index="0">
          ${badgeHTML}
        </div>
        <div class="item-details">
            <div class="item-header">
                <h2 class="item-name">${item.name}</h2>
            <div class="item-rating">
              ${'<i class="fas fa-star"></i>'.repeat(Math.floor(item.rating.avg_rating))}
                    ${item.rating.avg_rating % 1 >= 0.5 ? '<i class="fas fa-star-half-alt"></i>' : ''}
                    ${'<i class="far fa-star"></i>'.repeat(5 - Math.ceil(item.rating.avg_rating))}
                    <span>${item.rating.count_review}</span>
            </div>
          </div>
          <div class="item-category">${item.category.name}</div>
          <div class="item-category">${item.brand.name}</div>
          <div class="item-description">${item.description}</div>
          <div class="item-footer">
            <div class="item-price">
              ${item.price.toLocaleString('ru-RU')} ₽
            </div>
            <div class="item-actions">
            <a href="/brand/${item.brand.name}/product/${item.product_id}" class="category-link">
                    Подробнее <i class="fas fa-arrow-right"></i>
                </a>
                <div class="action-btn share-btn" data-product-id="${item.product_id}" data-brand-name="${item.brand.name}">
                    <i class="fas fa-share-alt"></i>
                </div>
                <div class="action-btn remove-btn" data-id="${item.product_id}" data-url="favorites">
                    <i class="fas fa-trash-alt"></i>
                </div>
            </div>
          </div>
        </div>
      `;
        return div;
    }

    // Удаление одного товара
    document.querySelectorAll('.remove-btn').forEach(button => {
        button.addEventListener('click', function() {
            const itemDiv = this.closest('.favorites-item');
            itemDiv.style.opacity = '0';
            setTimeout(async () => {
                itemDiv.style.display = 'none';
                const productId = button.dataset.id;
                try {
                    const res = await fetch(`/api/v1/favorites/${productId}`, {
                        method: 'DELETE',
                    });
                    if (!res.ok) throw new Error();
                } catch {
                    const error = new Error('Не удалось обновить избранное');
                    console.error(error);
                }
            }, 300);
        });
    });

    // Загрузка данных и рендер карточек
    fetch('/api/v1/favorites', {
        method: 'GET',
    })
        .then(res => {
            if (!res.ok) throw new Error('Ошибка загрузки избранного');
            return res.json();
        })
        .then(data => {
            console.log(data)
            listContainer.innerHTML = '';
            if (!Array.isArray(data) || data.length === 0) {
                listContainer.innerHTML = `
            <div class="empty-favorites">
              <div class="empty-icon"><i class="far fa-heart"></i></div>
              <h2 class="empty-title">Ваше избранное пусто</h2>
              <p class="empty-text">
                Добавляйте понравившиеся товары в избранное, чтобы не потерять их.
                Здесь будут храниться все товары, которые вам интересны.
              </p>
              <a href="/" class="btn btn-primary">
                <i class="fas fa-shopping-bag"></i> Перейти к покупкам
              </a>
            </div>
          `;
                return;
            }
            data.forEach(item => {
                const card = createFavoritesItem(item);
                listContainer.appendChild(card);
            });
            attachItemEventListeners();
        })
        .catch(err => {
            console.error(err);
            listContainer.innerHTML = `<p style="padding:20px">Ошибка загрузки избранного.</p>`;
        });

    // Кнопка "Поделиться списком"
    // document.getElementById('share-all-btn').addEventListener('click', function () {
    //     alert('Список избранных товаров отправлен!');
    // });

    // Кнопка "Очистить избранное"
    document.getElementById('clear-all-btn').addEventListener('click', function () {
        if (confirm('Вы уверены, что хотите очистить всё избранное?')) {
            setTimeout( async () => {
                try {
                    const res = await fetch(`/api/v1/favorites`, {
                        method: 'DELETE',
                    });
                    if (!res.ok) throw new Error();
                    listContainer.innerHTML = `
          <div class="empty-favorites">
            <div class="empty-icon"><i class="far fa-heart"></i></div>
            <h2 class="empty-title">Ваше избранное пусто</h2>
            <p class="empty-text">
              Добавляйте понравившиеся товары в избранное, чтобы не потерять их.
              Здесь будут храниться все товары, которые вам интересны.
            </p>
            <a href="/js/public" class="btn btn-primary">
              <i class="fas fa-shopping-bag"></i> Перейти на главную
            </a>
          </div>
        `;
                } catch {
                    alert('Не удалось очистить избранное.\n Повторите попытку позже.');
                    const error = new Error('Не удалось очистить избранное');
                    console.error(error);
                }
            }, 300);
        }
    });
});