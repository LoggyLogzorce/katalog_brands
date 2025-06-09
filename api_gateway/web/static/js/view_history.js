document.addEventListener('DOMContentLoaded', function () {
    const listContainer = document.getElementById('view-history-list');

    // Создаёт DOM-элемент для карточки товара
    function createViewHistoryItem(item) {
        const div = document.createElement('div');
        div.className = 'view-history-item';

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
              
            </div>
          </div>
          <div class="item-category">${item.category.name}</div>
          <div class="item-category">${item.brand_id}</div>
          <div class="item-description">${item.description}</div>
          <div class="item-footer">
            <div class="item-price">
              ${item.price.toLocaleString('ru-RU')} ₽
            </div>
            <div class="item-actions">
                <div class="action-btn share-btn" data-product-id="${item.product_id}" data-brand-id="${item.brand_id}">
                    <i class="fas fa-share-alt"></i>
                </div>
                <div class="action-btn remove-btn" data-id="${item.product_id}" data-url="view-history">
                    <i class="fas fa-trash-alt"></i>
                </div>
            </div>
          </div>
        </div>
      `;
        return div;
    }

    // Загрузка данных и рендер карточек
    fetch('/api/v1/view-history', {
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
            <div class="empty-view-history">
            <div class="empty-icon"><i class="far fa-heart"></i></div>
            <h2 class="empty-title">Ваша история просмотра пуста</h2>
            <p class="empty-text">
              Здесь будут храниться все товары, которые вы просматривали.
            </p>
            <a href="/" class="btn btn-primary">
              <i class="fas fa-shopping-bag"></i> Перейти на главную
            </a>
          </div>
          `;
                return;
            }
            data.forEach(item => {
                const card = createViewHistoryItem(item);
                listContainer.appendChild(card);
            });
            attachItemEventListeners();
        })
        .catch(err => {
            console.error(err);
            listContainer.innerHTML = `<p style="padding:20px">Ошибка загрузки избранного.</p>`;
        });

    // Кнопка "Очистить избранное"
    document.getElementById('clear-all-btn').addEventListener('click', function () {
        if (confirm('Вы уверены, что хотите очистить историю просмотра?')) {
            setTimeout(async () => {
                try {
                    const res = await fetch(`/api/v1/view-history`, {
                        method: 'DELETE',
                    });
                    if (!res.ok) throw new Error();
                    listContainer.innerHTML = `
          <div class="empty-view-history">
            <div class="empty-icon"><i class="far fa-heart"></i></div>
            <h2 class="empty-title">Ваша история просмотра пуста</h2>
            <p class="empty-text">
              Здесь будут храниться все товары, которые вы просматривали.
            </p>
            <a href="/" class="btn btn-primary">
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

            // Обновляем счётчик в шапке, если он есть:
            const countEl = document.querySelector('.favorites-count');
            if (countEl) countEl.textContent = '0';
            // Дополнительно: отправить запрос на очистку списка на сервере
        }
    });
});