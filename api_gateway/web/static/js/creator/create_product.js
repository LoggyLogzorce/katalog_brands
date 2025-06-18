document.addEventListener('DOMContentLoaded', () => {
    const brandName = window.location.pathname.split('/')[3]; // /creator/brand/{name}
    const openBtn = document.getElementById('createProductBtn');
    const modal = document.getElementById('createProductModal');
    const closeBtn = modal.querySelector('.close-modal');
    const cancelBtn = document.getElementById('cancelCreateProduct');
    const form = document.getElementById('createProductForm');
    const categorySelect = document.getElementById('productCategory');
    const imagesInput = document.getElementById('productImages');
    const previewContainer = document.getElementById('imagesPreviewContainer');

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

    // Открыть модалку
    openBtn.addEventListener('click', () => {
        modal.classList.add('active');
    });

    // Закрыть модалку
    [closeBtn, cancelBtn].forEach(btn =>
        btn.addEventListener('click', () => modal.classList.remove('active'))
    );
    modal.addEventListener('click', e => {
        if (e.target === modal) modal.classList.remove('active');
    });

    document.addEventListener('keydown', e => {
        if (e.key === 'Escape' && modal.classList.contains('active')) {
            modal.classList.remove('active');
        }
    });

    // Загрузить список категорий
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

    // Отправка формы
    form.addEventListener('submit', async e => {
        e.preventDefault();
        const fd = new FormData(form);

        try {
            const res = await fetch(`/api/v1/creator/brand/${encodeURIComponent(brandName)}/create-product`, {
                method: 'POST',
                body: fd
            });
            if (!res.ok) throw new Error('Ошибка создания');

            // Закрыть модалку
            modal.classList.remove('active');
            form.reset();

            window.location.reload();
        } catch (err) {
            alert(err.message);
        }
    });
});
