document.addEventListener('DOMContentLoaded', function () {
    const parts = window.location.pathname.split('/');
    const count = parts[parts.length - 1] || parts[parts.length - 2];
    let param = 'all';
    if (count !== 'brands') {
        param = 3
    }
    const grid = document.querySelector('.brands-grid');
    grid.innerHTML = '<p class="loading">Загрузка брендов…</p>';

    // Запрос на получение категорий
    fetch(`/api/v1/brands?count=${param}`, {
        method: 'GET'
    })
        .then(res => {
            if (!res.ok) throw new Error('Не удалось загрузить категории');
            return res.json();
        })
        .then(brands => {
            grid.innerHTML = '';
            if (!Array.isArray(brands) || brands.length === 0) {
                grid.innerHTML = '<p>Категории не найдены.</p>';
                return;
            }
            brands.forEach(b => {
                const card = document.createElement('div');
                card.className = 'brand-card';
                card.innerHTML = `<div class="brand-image">
                <img src="/static/${b.logo_url}" alt="${b.name}">
                <div class="overlay"></div>
            </div>
            <div class="brand-info">
                <div>
                    <h3 class="brand-name">${b.name}</h3>
                    <div class="brand-products">${b.product_count} товаров</div>
                </div>
                <a href="/brand/${b.name}" class="brand-link">
                    Смотреть все <i class="fas fa-arrow-right"></i>
                </a>
            </div>`;
                grid.appendChild(card);
            });

            // Наведение и клик
            const brandCards = grid.querySelectorAll('.brand-card');
            brandCards.forEach(card => {
                const link = card.querySelector('.brand-link');
                card.addEventListener('mouseenter', () => link.style.gap = '12px');
                card.addEventListener('mouseleave', () => link.style.gap = '8px');
                card.addEventListener('click', e => {
                    if (!e.target.closest('.brand-link')) {
                        window.location.href = link.href;
                    }
                });
            });
        })
        .catch(err => {
            console.error(err);
            grid.innerHTML = '<p class="error">Ошибка загрузки брендов.</p>';
        });
})