function attachItemEventListeners() {
    // Переключение фотографий по клику на изображение
    document.querySelectorAll('.item-image').forEach(img => {
        img.addEventListener('click', function() {
            const images = JSON.parse(this.getAttribute('data-images'));
            let idx = parseInt(this.getAttribute('data-current-index'), 10);
            if (images.length <= 1) return;
            idx = (idx + 1) % images.length;
            this.src = '/static/' + images[idx];
            this.setAttribute('data-current-index', idx);
        });
    });

    // Удаление одного товара
    document.querySelectorAll('.remove-btn').forEach(button => {
        button.addEventListener('click', function() {
            const url = button.dataset.url;
            const cls = url === 'view-history' ? '.view-history-item' : '.favorites-item';
            const itemDiv = this.closest(cls);
            itemDiv.style.opacity = '0';
            setTimeout(async () => {
                itemDiv.style.display = 'none';
                const productId = button.dataset.id;
                try {
                    const res = await fetch(`/api/v1/${url}/${productId}`, {
                        method: 'DELETE',
                    });
                    if (!res.ok) throw new Error();
                } catch {
                    const error = new Error('Не удалось обновить избранное');
                    console.error(error);
                }
            }, 300);
        });
    });

    // Поделиться одним товаром
    document.querySelectorAll('.share-btn').forEach(button => {
        button.addEventListener('click', function() {
            const productId = this.getAttribute('data-product-id');
            const brandId = this.getAttribute('data-brand-name');
            const productURL = `${window.location.origin}/brand/${brandId}/product/${productId}`;
            navigator.clipboard.writeText(productURL)
                .then(() => {
                    alert('Ссылка на товар скопирована:\n' + productURL);
                })
                .catch(() => {
                    alert('Не удалось скопировать ссылку');
                });
        });
    });
}

// Проверка, входит ли дата в последние 2 недели
function isNewProduct(createdAt) {
    const twoWeeksAgo = new Date();
    twoWeeksAgo.setDate(twoWeeksAgo.getDate() - 14);
    return new Date(createdAt) >= twoWeeksAgo;
}