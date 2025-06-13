document.addEventListener('DOMContentLoaded', function () {
    const grid = document.getElementById('products-list');
    grid.innerHTML = '<p class="loading">Загрузка товаров…</p>';

    function createProductCard(item) {
        const badgeHTML = isNewProduct(item.created_at)
            ? '<span class="badge">Новинка</span>'
            : '';

        const div = document.createElement('div');
        div.className = 'product-card';
        div.innerHTML = `
            ${badgeHTML}
            <img src="/static/${item.product_urls[0].url}"
                 alt="${item.name}" class="product-image">
            <div class="product-info">
                <h3 class="product-name">${item.name}</h3>
                <div class="product-rating">
                    ${'<i class="fas fa-star"></i>'.repeat(Math.floor(item.rating.avg_rating))}
                    ${item.rating.avg_rating % 1 >= 0.5 ? '<i class="fas fa-star-half-alt"></i>' : ''}
                    ${'<i class="far fa-star"></i>'.repeat(5 - Math.ceil(item.rating.avg_rating))}
                    <span>${item.rating.count_review}</span>
                </div>
                <div class="product-price">${item.price.toLocaleString('ru-RU')} ₽</div>
                <div class="product-actions">
                    <div class="action-btn favorite-btn" data-id="${item.product_id}">
                        <i class="${item.is_favorite ? 'fas' : 'far'} fa-heart" style="color: #FFB6C1;"></i>
                    </div>
                </div>
            </div>`;
        return div;
    }

    function renderProducts(list) {
        grid.innerHTML = '';
        if (!list.length) {
            grid.innerHTML = '<p>Товары не найдены.</p>';
            return;
        }
        list.forEach(item => {
            grid.appendChild(createProductCard(item));
        });
        attachFavoriteHandlers();
    }

    function applySort() {
        const mode = sortSelect.value;
        let sorted = [...products];
        if (mode === 'price-asc') {
            sorted.sort((a, b) => a.price - b.price);
        } else if (mode === 'price-desc') {
            sorted.sort((a, b) => b.price - a.price);
        } else if (mode === 'newest') {
            sorted.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
        } else if (mode === 'rating') {
            sorted.sort((a, b) => (b.rating.avg_rating || 0) - (a.rating.avg_rating || 0));
        }
        renderProducts(sorted);
    }

    const parts = window.location.pathname.split('/');
    const count = parts[parts.length - 1] || parts[parts.length - 2];
    let param = 'all';
    if (count !== '/') {
        param = 4
    }

    fetch(`/api/v1/products/approved?count=${param}`, {method: 'GET'})
        .then(res => {
            if (!res.ok) throw new Error('Ошибка загрузки товаров');
            return res.json();
        })
        .then(data => {
            grid.innerHTML = '';
            if (!Array.isArray(data) || data.length === 0) {
                grid.innerHTML = '<p>Товары в этой категории не найдены.</p>';
                return;
            }
            renderProducts(data);
        })
        .catch(err => {
            console.error(err);
            grid.innerHTML = '<p class="error">Не удалось загрузить товары.</p>';
        });
});
