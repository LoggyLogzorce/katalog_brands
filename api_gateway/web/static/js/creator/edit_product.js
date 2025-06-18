document.addEventListener('DOMContentLoaded', () => {
    const brandName = window.location.pathname.split('/')[3];
    const modal = document.getElementById('editProductModal');
    const closeBtns = modal.querySelectorAll('.close-modal, #cancelEditProduct');
    const form = document.getElementById('editProductForm');
    const nameInput = document.getElementById('newProductName');
    const descInput = document.getElementById('newProductDesc');
    const priceInput = document.getElementById('newProductPrice');
    const categorySelect = document.getElementById('newProductCategory');
    const imagesInput = document.getElementById('newProductImages');
    const previewContainer = document.getElementById('newImagesPreviewContainer');
    let currentProductId = null;

    // закрытие модалки
    closeBtns.forEach(b => b.addEventListener('click', () => modal.classList.remove('active')));
    modal.addEventListener('click', e => { if (e.target === modal) modal.classList.remove('active'); });
    document.addEventListener('keydown', e => { if (e.key==='Escape') modal.classList.remove('active'); });

    // делегируем клик по edit-кнопкам внутри grid
    document.querySelector('.products-grid').addEventListener('click', async e => {
        const btn = e.target.closest('.edit-btn');
        if (!btn) return;
        e.preventDefault();

        // получаем ID и данные товара
        currentProductId = btn.dataset.id;
        const res = await fetch(`/api/v1/brand/${encodeURIComponent(brandName)}/product/${currentProductId}?status=creator`);
        const { product: p } = await res.json();

        // заполняем форму
        nameInput.value = p.name;
        descInput.value = p.description;
        priceInput.value = p.price;
        categorySelect.value = p.category_id;

        // превью старых картинок
        previewContainer.innerHTML = '';
        p.product_urls.forEach(u => {
            const div = document.createElement('div');
            div.className = 'image-preview';
            div.innerHTML = `
                <img src="/static/${u.url}" alt="">
                <button class="remove-old-image" data-url="${u.url}">&times;</button>
            `;
            previewContainer.append(div);
        });

        // открываем модалку
        modal.classList.add('active');
    });

    // Предпросмотр изображений при выборе файлов
    imagesInput.addEventListener('change', function(e) {
        previewContainer.innerHTML = '';

        const files = this.files;
        if (files.length > 0) {
            for (let i = 0; i < files.length; i++) {
                const file = files[i];
                if (file.type.match('image.*')) {
                    const reader = new FileReader();

                    reader.onload = function(e) {
                        const preview = document.createElement('div');
                        preview.className = 'image-preview fade-in';
                        preview.innerHTML = `
                                <img src="${e.target.result}" alt="Предпросмотр">
                                <button class="remove-image" data-index="${i}">&times;</button>
                            `;
                        previewContainer.appendChild(preview);

                        // Добавляем обработчик для удаления изображения
                        const removeBtn = preview.querySelector('.remove-image');
                        removeBtn.addEventListener('click', function() {
                            // Удаляем файл из input
                            const filesArray = Array.from(imagesInput.files);
                            filesArray.splice(parseInt(this.dataset.index), 1);

                            // Создаем новый FileList
                            const dataTransfer = new DataTransfer();
                            filesArray.forEach(file => dataTransfer.items.add(file));
                            imagesInput.files = dataTransfer.files;

                            // Удаляем превью
                            preview.remove();
                        });
                    }

                    reader.readAsDataURL(file);
                }
            }
        }
    });

    fetch(`/api/v1/categories`, {
        method: 'GET'
    })
        .then(r => r.json())
        .then(list => {
            list.forEach(cat => {
                const opt = document.createElement('option');
                opt.value = cat.category_id;
                opt.textContent = cat.name;
                categorySelect.append(opt);
            });
        });


    // отправка изменений
    form.addEventListener('submit', async e => {
        e.preventDefault();
        const fd = new FormData(form);
        const res = await fetch(
            `/api/v1/creator/brand/${encodeURIComponent(brandName)}/product/${currentProductId}`,
            { method: 'PUT', body: fd }
        );
        if (!res.ok) return alert('Ошибка');
        modal.classList.remove('active');
        window.location.reload();
    });
});
