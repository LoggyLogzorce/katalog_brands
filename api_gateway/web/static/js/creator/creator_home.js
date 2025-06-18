document.addEventListener('DOMContentLoaded', () => {
    const stats = {
        total_brands: document.querySelector('.stats-container .stat-card:nth-child(1) .stat-value'),
        total_products: document.querySelector('.stats-container .stat-card:nth-child(2) .stat-value'),
        overall_rating: document.querySelector('.stats-container .stat-card:nth-child(3) .stat-value'),
        total_views: document.querySelector('.stats-container .stat-card:nth-child(4) .stat-value'),
    };
    const grid = document.querySelector('.brands-grid');

    function formatNumber(n) {
        return n.toLocaleString('ru-RU');
    }

    fetch('/api/v1/creator/brands')
        .then(res => res.json())
        .then(data => {
            const brands = data.brands;
            // Статистика
            stats.total_brands.textContent = brands.length;
            stats.total_products.textContent = brands.reduce((sum, b) => sum + b.products_count, 0);
            stats.overall_rating.textContent = (brands.reduce((s, b) => s + b.avg_rating, 0) / brands.length || 0).toFixed(1);
            stats.total_views.textContent = formatNumber(brands.reduce((s, b) => s + b.views, 0));

            // Очистим и заполним сетку
            grid.innerHTML = '';
            brands.forEach(b => {
                const card = document.createElement('div');
                card.className = 'brand-card';
                card.innerHTML = `
          <div class="brand-image">
            <img src="/static/${b.logo_url}" alt="${b.name}">
            <div class="overlay"></div>
          </div>
          <div class="brand-info">
            <div>
              <h3 class="brand-name">${b.name}</h3>
              <div class="brand-products">${b.products_count} ${b.products_count === 1 ? 'товар' : b.products_count < 5 ? 'товара' : 'товаров'}</div>
            </div>
            <div class="brand-stats" style="display:flex;gap:8px;font-size:.9rem;color:#777">
              <span>★ ${b.avg_rating.toFixed(1)}</span>
              <span>👁 ${formatNumber(b.views)}</span>
            </div>
            <a href="/creator/brand/${encodeURIComponent(b.name)}" class="brand-link">
              Управлять брендом <i class="fas fa-arrow-right"></i>
            </a>
          </div>
        `;
                grid.append(card);
            });
        })
        .catch(err => {
            console.error('Не удалось загрузить мои бренды:', err);
            grid.innerHTML = '<p style="padding:20px">Ошибка загрузки брендов.</p>';
        });

    document.querySelector('.create-brand-btn').addEventListener('click', () => {
        // Перейти на форму создания
        window.location.href = '/brand/create';
    });
});