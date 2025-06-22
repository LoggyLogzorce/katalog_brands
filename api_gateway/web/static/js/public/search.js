document.addEventListener('DOMContentLoaded', () => {
    const input = document.querySelector('.search-input');
    const btn = document.querySelector('.search-btn');
    const heroDesc = document.querySelector('.brand-description');
    const noResSec = document.querySelector('.no-results');
    const prodSec = document.querySelector('.products');
    const brandSec = document.querySelector('.brands');
    const catSec = document.querySelector('.categories');
    const prodGrid = document.querySelector('.products-grid');
    const brandGrid = document.querySelector('.brands-grid');
    const catGrid = document.querySelector('.categories-grid');

    async function doSearch(q) {
        input.valueOf().value = q;
        const res = await fetch(`/api/v1/search?q=${encodeURIComponent(q)}&size=10&from=0`);
        if (!res.ok) {
            prodSec.style.display = 'none';
            brandSec.style.display = 'none';
            catSec.style.display = 'none';
            return;
        }
        const {products, brands, categories} = await res.json();
        // products.items, brands.items, categories.items
        if (products.total === 0) {
            prodSec.style.display = 'none';
        } else {
            renderProducts(products.items);
            attachFavoriteHandlers()
            prodSec.style.display = 'block';
        }

        if (brands.total === 0) {
            brandSec.style.display = 'none';
        } else {
            renderBrands(brands.items);
            brandSec.style.display = 'block';
        }

        if (categories.total === 0) {
            catSec.style.display = 'none';
        } else {
            renderCategories(categories.items);
            catSec.style.display = 'block';
        }
        heroDesc.textContent =
            `Найдено ${products.total} товара, ${brands.total} брендов и ${categories.total} категорий`;
    }

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
        prodGrid.innerHTML = '';

        products.forEach(product => {
            prodGrid.innerHTML += createProductCard(product);
        });
    }

    function renderBrands(items) {
        brandGrid.innerHTML = items.map(b => `
      <div class="brand-card">
        <div class="brand-image">
          <img src="/static/${b.logo_url}" alt="${b.name}">
          <div class="overlay"></div>
        </div>
        <div class="brand-info">
          <h3 class="brand-name">${b.name}</h3>
          <p class="brand-products">${b.product_count} товара</p>
          <a href="/brand/${b.name}" class="brand-link">
            Посмотреть товары <i class="fas fa-arrow-right"></i>
          </a>
        </div>
      </div>
    `).join('');
    }

    function renderCategories(items) {
        catGrid.innerHTML = items.map(c => `
      <div class="category-card">
        <div class="category-image">
          <img src="/static/${c.photo}" alt="${c.name}">
          <div class="overlay"></div>
        </div>
        <div class="category-info">
          <h3 class="category-name">${c.name}</h3>
          <p class="category-products">${c.product_count} товаров</p>
          <a href="/category/${c.category_id}" class="category-link">
            Посмотреть товары <i class="fas fa-arrow-right"></i>
          </a>
        </div>
      </div>
    `).join('');
    }

    btn.addEventListener('click', () => doSearch(input.value.trim()));
    input.addEventListener('keydown', e => {
        if (e.key === 'Enter') doSearch(input.value.trim());
    });

    // инициалный поиск на странице
    const params = new URLSearchParams(window.location.search);
    const q = params.get('q');
    doSearch(q);
});
