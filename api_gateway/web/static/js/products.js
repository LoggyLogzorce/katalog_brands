document.addEventListener('DOMContentLoaded', function() {
    const parts = window.location.pathname.split('/');
    const productID = parts[parts.length - 1] || parts[parts.length - 2];
    const sortSelect = document.getElementById('sort');
    let products = [];

    let param = 'all'
    if (productID !== 'products') {
        param = 4
    }

    fetch(`/api/v1/products/approved?count=${param}`, {
        method: 'GET',
    })
        .then(response => response.json())
        .then(products => {
            initPage(products);
        })
        .catch(error => {
            console.error(error);
        });

    // Функция для создания карточки товара
    function createProductCard(product) {
        const mainImage = product.product_urls.length > 0
            ? product.product_urls[0].url
            : 'https://via.placeholder.com/300x200?text=No+Image';

        const badgeHTML = isNewProduct(product.created_at)
            ? '<span class="badge">Новинка</span>'
            : '';

        return `
            <div class="product-card" data-product-id="${product.product_id}">
                ${badgeHTML}
                <img src="/static/${mainImage}" alt="${product.name}" class="product-image">
                <div class="product-info">
                    <span class="product-category">${product.category.name}</span>
                    <h3 class="product-name">${product.name}</h3>
                    <div class="product-rating">
                        ${renderRatingStars(product.rating.avg_rating)}
                        <span>(${product.rating.count_review})</span>
                    </div>
                    <div class="product-price">${formatPrice(product.price)}</div>
                    <div class="product-actions">
                        <div class="action-btn favorite-btn" data-id="${product.product_id}">
                            <i class="${product.is_favorite ? 'fas' : 'far'} fa-heart" style="color: #FFB6C1;"></i>
                        </div>
                    </div>
                    <a href="/brand/${product.brand.name}/product/${product.product_id}" class="category-link">
                    Подробнее <i class="fas fa-arrow-right"></i>
                </a>
                </div>
            </div>
        `;
    }

    // Функция для рендеринга всех товаров
    function renderProducts(products) {
        const productsGrid = document.querySelector('.products-grid');
        productsGrid.innerHTML = '';

        products.forEach(product => {
            productsGrid.innerHTML += createProductCard(product);
        });
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

    // Инициализация страницы
    function initPage(data) {
        renderProducts(data);
        products = data;
        attachFavoriteHandlers()
    }

    sortSelect.addEventListener('change', applySort);
});