document.addEventListener('DOMContentLoaded', () => {
    const parts = window.location.pathname.split('/');
    const brandName = parts[parts.indexOf('creator') + 2]; // /creator/brand/{name}

    const logoImg = document.querySelector('.brand-logo');
    const nameEl = document.querySelector('.brand-name');
    const statusEl = document.querySelector('.moderation-status');
    const descEl = document.querySelector('.brand-description');
    const metaItems = document.querySelectorAll('.brand-meta .meta-value');
    const productsHeader = document.querySelector('.products-header h3');
    const grid = document.querySelector('.products-grid');

    function formatDate(iso) {
        const d = new Date(iso);
        return d.toLocaleDateString('ru-RU');
    }

    function formatPrice(p) {
        return p.toLocaleString('ru-RU') + ' ₽';
    }

    fetch(`/api/v1/creator/brand/${encodeURIComponent(brandName)}`)
        .then(res => res.json())
        .then(data => {
            const b = data.brand;
            // Левый блок
            logoImg.src = `/static/${b.logo_url}`;
            nameEl.textContent = b.name;
            // статус, если есть b.status
            if (b.status === 'approved') {
                statusEl.className = 'moderation-status status-approved';
                statusEl.innerHTML = '<i class="fas fa-check-circle"></i> Одобрено модератором';
            } else if (b.status === 'pending') {
                statusEl.className = 'moderation-status status-pending';
                statusEl.innerHTML = '<i class="fas fa-hourglass-half"></i> На модерации';
            } else if (b.status === 'rejected') {
                statusEl.className = 'moderation-status status-rejected';
                statusEl.innerHTML = '<i class="fas fa-times-circle"></i> Отклонено модератором';
            }
            descEl.textContent = b.description || '';

            metaItems[0].textContent = formatDate(b.created_at);          // Дата создания
            metaItems[1].textContent = b.products_count;                  // Товаров
            metaItems[2].textContent = b.views.toLocaleString('ru-RU');   // Просмотры
            metaItems[3].innerHTML = `${b.avg_rating.toFixed(1)} <i class="fas fa-star" style="color:#FFD700"></i>`;

            // Правый блок: заголовок и кнопка добавления
            productsHeader.textContent = `Товары бренда (${b.products_count})`;

            // Товары
            grid.innerHTML = '';
            data.product.forEach(p => {
                let status = '';
                let statusText = '';
                if (p.status === 'approved') {
                    status = 'fa-check';
                    statusText = 'Одобрено';
                } else if (p.status === 'pending') {
                    status = 'fa-pencil-alt';
                    statusText = 'На модерации';
                } else if (p.status === 'rejected') {
                    status = 'fa-times';
                    statusText = 'Отклонено';
                }
                const div = document.createElement('div');
                div.className = 'product-card';
                div.innerHTML = `
          <div class="product-status status-${p.status}">
            <i class="fas ${status}"></i>
            ${statusText}
          </div>
          <img src="/static/${p.product_urls[0].url}" class="product-image" alt="${p.name}">
          <div class="product-info">
            <h3 class="product-name">${p.name}</h3>
            <div class="product-price">${formatPrice(p.price)}</div>
            <div class="product-actions">
              <button data-id="${p.product_id}" class="action-btn edit-btn">
                <i class="fas fa-edit"></i> Редактировать
              </button>
              <button data-id="${p.product_id}" class="action-btn delete-btn">
                <i class="fas fa-trash"></i> Удалить
              </button>
            </div>
          </div>
        `;
                grid.appendChild(div);
            });

            // Удаление товара
            document.querySelectorAll('.delete-btn').forEach(btn => {
                btn.addEventListener('click', async e => {
                    e.preventDefault();
                    const id = btn.dataset.id;
                    if (!confirm('Удалить этот товар?')) return;
                    const res = await fetch(`/api/v1/creator/brand/${encodeURIComponent(b.name)}/product/${id}`, {
                        method: 'DELETE'
                    });
                    if (res.ok) btn.closest('.product-card').remove();
                    else alert('Ошибка при удалении');
                });
            });
        })
        .catch(err => {
            console.error('Ошибка загрузки данных:', err);
        });

    const openBtn = document.getElementById('editBrandBtn');
    const modalEditBrand = document.getElementById('editBrandModal');
    const closeBtn = document.querySelector('.close-modal');
    const cancelBtn = document.getElementById('cancelEditBrand');
    const form = document.getElementById('editBrandForm');
    const logoInput = document.getElementById('brandLogoInput');
    const logoPreview = document.getElementById('logoPreview');

    // Открыть модалку и заполнить текущими данными
    openBtn.addEventListener('click', () => {
        document.getElementById('brandNameInput').value = document.querySelector('.brand-name').textContent;
        document.getElementById('brandDescInput').value = document.querySelector('.brand-description').textContent;

        // Установить превью логотипа
        logoPreview.src = document.querySelector('.brand-logo').src;
        logoPreview.style.display = 'block';

        modalEditBrand.classList.add('active');
    });

    // Закрыть модалку
    function closeModal() {
        modalEditBrand.classList.remove('active');
    }

    // События закрытия
    closeBtn.addEventListener('click', closeModal);
    cancelBtn.addEventListener('click', closeModal);

    // Закрыть при клике вне модального окна
    modalEditBrand.addEventListener('click', (e) => {
        if (e.target === modalEditBrand) {
            closeModal();
        }
    });

    // Закрыть по клавише Esc
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && modalEditBrand.classList.contains('active')) {
            closeModal();
        }
    });

    // Предпросмотр логотипа при выборе файла
    logoInput.addEventListener('change', function (e) {
        const file = this.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = function (e) {
                logoPreview.src = e.target.result;
                logoPreview.style.display = 'block';
            }
            reader.readAsDataURL(file);
        }
    });

    // Отправка формы
    form.addEventListener('submit', async e => {
        e.preventDefault();

        const url = `/api/v1/creator/brand/${encodeURIComponent(brandName)}/edit`;
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
                method: 'PUT',
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

            alert('Изменения бренда успешно сохранены!');

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
