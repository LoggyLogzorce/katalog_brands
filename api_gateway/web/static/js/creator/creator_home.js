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

    const openBtn = document.getElementById('createBrandBtn');
    const modal = document.getElementById('createBrandModal');
    const closeBtn = document.querySelector('.close-modal');
    const cancelBtn = document.getElementById('cancelCreateBrand');
    const form = document.getElementById('createBrandForm');
    const logoInput = document.getElementById('brandLogoInput');
    const logoPreview = document.getElementById('logoPreview');

    // Открыть модалку и заполнить текущими данными
    openBtn.addEventListener('click', () => {
        modal.classList.add('active');
    });

    // Закрыть модалку
    function closeModal() {
        modal.classList.remove('active');
    }

    // События закрытия
    closeBtn.addEventListener('click', closeModal);
    cancelBtn.addEventListener('click', closeModal);

    // Закрыть при клике вне модального окна
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            closeModal();
        }
    });

    // Закрыть по клавише Esc
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && modal.classList.contains('active')) {
            closeModal();
        }
    });

    // Предпросмотр логотипа при выборе файла
    logoInput.addEventListener('change', function(e) {
        const file = this.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = function(e) {
                logoPreview.src = e.target.result;
                logoPreview.style.display = 'block';
            }
            reader.readAsDataURL(file);
        }
    });

    // Отправка формы
    form.addEventListener('submit', async e => {
        e.preventDefault();

        const url = `/api/v1/creator/brand/create`;
        const formData = new FormData(form);
        for (let [key, value] of formData.entries()) {
            console.log(key, value);
        }
        try {
            // Отключаем кнопку «Сохранить» на время запроса
            const submitBtn = form.querySelector('button[type="submit"]');
            submitBtn.disabled = true;
            submitBtn.textContent = 'Сохраняем...';

            const response = await fetch(url, {
                method: 'POST',
                body: formData,
                credentials: 'include', // если нужны куки
            });

            if (!response.ok) {
                const err = await response.json().catch(() => ({}));
                throw new Error(err.error || 'Ошибка сохранения');
            }

            const updated = await response.json();
            // Закрываем модалку
            closeModal();

            alert('Бренд успешно создан!');

            window.location.href = `/creator/brand/${encodeURIComponent(updated.name)}`;
        } catch (err) {
            console.error(err);
            alert(err.message);
        } finally {
            const submitBtn = form.querySelector('button[type="submit"]');
            submitBtn.disabled = false;
            submitBtn.innerHTML = '<i class="fas fa-save"></i> Сохранить';
        }
    });
});