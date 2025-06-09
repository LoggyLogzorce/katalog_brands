document.addEventListener('DOMContentLoaded', function () {
    // Поиск по категориям
    const searchInput = document.querySelector('.search-bar input');
    const searchButton = document.querySelector('.search-bar button');

    searchButton.addEventListener('click', function () {
        if (searchInput.value.trim() !== '') {
            alert(`Поиск: "${searchInput.value}"\nРезультаты будут показаны на следующей странице.`);
            // window.location.href = `/search?query=${encodeURIComponent(searchInput.value)}`;
        } else {
            searchInput.focus();
        }
    });

    searchInput.addEventListener('keypress', function (e) {
        if (e.key === 'Enter') searchButton.click();
    });

    // Контейнер для категорий
    const grid = document.querySelector('.categories-grid');
    grid.innerHTML = '<p class="loading">Загрузка категорий…</p>';

    // Запрос на получение категорий
    fetch('/api/v1/categories', {
        method: 'GET'
    })
        .then(res => {
            if (!res.ok) throw new Error('Не удалось загрузить категории');
            return res.json();
        })
        .then(categories => {
            grid.innerHTML = '';
            if (!Array.isArray(categories) || categories.length === 0) {
                grid.innerHTML = '<p>Категории не найдены.</p>';
                return;
            }
            categories.forEach(cat => {
                const card = document.createElement('div');
                card.className = 'category-card';
                card.innerHTML = `<div class="category-image">
                <img src="/static/${cat.photo}" alt="${cat.name}">
                <div class="overlay"></div>
            </div>
            <div class="category-info">
                <div>
                    <h3 class="category-name">${cat.name}</h3>
                    <div class="category-products">${cat.product_count} товаров</div>
                </div>
                <a href="/category/${cat.category_id}" class="category-link">
                    Смотреть все <i class="fas fa-arrow-right"></i>
                </a>
            </div>`;
                grid.appendChild(card);
            });

            // Наведение и клик
            const categoryCards = grid.querySelectorAll('.category-card');
            categoryCards.forEach(card => {
                const link = card.querySelector('.category-link');
                card.addEventListener('mouseenter', () => link.style.gap = '12px');
                card.addEventListener('mouseleave', () => link.style.gap = '8px');
                card.addEventListener('click', e => {
                    if (!e.target.closest('.category-link')) {
                        window.location.href = link.href;
                    }
                });
            });
        })
        .catch(err => {
            console.error(err);
            grid.innerHTML = '<p class="error">Ошибка загрузки категорий.</p>';
        });
});

