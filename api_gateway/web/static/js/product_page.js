document.addEventListener('DOMContentLoaded', () => {
    // Получаем productId и brandName из URL, например /brand/Skin%20Care/product/7
    const parts = window.location.pathname.split('/');
    const brandName = parts[parts.indexOf('brand') + 1];
    const productId = parts[parts.indexOf('product') + 1];
    let products = [];

    // Селекторы
    const mainImg        = document.getElementById('mainImage');
    const thumbs         = document.querySelector('.thumbnails');
    const titleEl        = document.querySelector('.product-title');
    const descEl         = document.querySelector('.product-description');
    const priceEl        = document.querySelector('.product-price');
    const reviewsList    = document.querySelector('.reviews-list');
    const overallRating  = document.querySelector('.overall-rating');
    const starsLarge     = document.querySelector('.stars-large');
    const reviewCount1   = document.querySelector('.review-count-1');
    const reviewCount2   = document.querySelector('.review-count-2');
    const viewAllLink    = document.querySelector('.view-all');
    const favoriteBtn    = document.getElementById('favoriteBtn');

    // Форматирование
    function formatPrice(p) {
        return p.toLocaleString('ru-RU') + ' ₽';
    }

    // Обновить кнопку «избранное»
    function updateFavorite(isFav) {
        favoriteBtn.classList.toggle('favorited', isFav);
        favoriteBtn.innerHTML = `<i class="${isFav ? 'fas' : 'far'} fa-heart"></i> `
            + (isFav ? 'В избранном' : 'Добавить в избранное');
    }

    // Рендер отзывов
    function renderReviews(reviews) {
        reviewsList.innerHTML = '';
        const count = reviews.length;
        const avg = count
            ? (reviews.reduce((s,r)=>s+r.rating,0)/count).toFixed(1)
            : '0.0';

        // Обновляем общий блок рейтинга
        overallRating.textContent = avg;
        starsLarge.innerHTML      = renderRatingStars(avg);
        reviewCount1.textContent  = `${count} отзыв${count===1?'':count<5?'a':'ов'}`;
        reviewCount2.textContent  = reviewCount1.textContent;

        // Выводим каждый отзыв
        reviews.forEach(r => {
            const div = document.createElement('div');
            div.className = 'review-item';
            div.innerHTML = `
        <div class="review-header">
          <div class="reviewer-name">${r.user.name}</div>
          <div class="review-date">${new Date(r.created_at).toLocaleDateString('ru-RU')}</div>
        </div>
        <div class="review-rating">
          ${renderRatingStars(r.rating)}
        </div>
        <div class="review-text">${r.description}</div>
      `;
            reviewsList.append(div);
        });
    }

    // Рендер миниатюр
    function renderThumbnails(urls) {
        thumbs.innerHTML = '';
        urls.forEach((u,i) => {
            const t = document.createElement('div');
            t.className = 'thumbnail' + (i===0?' active':'');
            t.innerHTML = `<img src="/static/${u.url}" alt="">`;
            t.addEventListener('click', () => {
                mainImg.src = `/static/${u.url}`;
                thumbs.querySelectorAll('.thumbnail').forEach(x=>x.classList.remove('active'));
                t.classList.add('active');
            });
            thumbs.append(t);
        });
    }

    // Обработчик клика «избранное»
    favoriteBtn.addEventListener('click', async () => {
        const isFav = favoriteBtn.classList.contains('favorited');
        try {
            const res = await fetch(`/api/v1/favorites/${productId}`, {
                method: isFav ? 'DELETE' : 'POST'
            });
            if (!res.ok) throw '';
            updateFavorite(!isFav);
        } catch {
            alert('Не удалось обновить избранное');
        }
    });

    fetch(`/api/v1/brand/${encodeURIComponent(brandName)}/product/${productId}`, {
        method: 'GET'
    })
        .then(res => res.json())
        .then(data => {
            const p = data.product;
            // Заполняем шапку товара
            titleEl.textContent       = p.name;
            descEl.textContent        = p.description;
            priceEl.textContent       = formatPrice(p.price);
            viewAllLink.href          = `/brand/${encodeURIComponent(data.product.brand.name)}`;
            // Изображения
            mainImg.src               = `/static/${p.product_urls[0]?.url}`;
            renderThumbnails(p.product_urls);
            // Избранное
            updateFavorite(p.is_favorite);
            // Отзывы
            renderReviews(data.reviews);
        })
        .catch(err => {
            console.error('Ошибка загрузки товара:', err);
        });

    fetch(`/api/v1/brand/${brandName}`, {
        method: 'GET',
    })
        .then(response => response.json())
        .then(brandData => {
            initPage(brandData);
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
                    <a href="/brand/${brandName}/product/${product.product_id}" class="category-link">
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

    // Инициализация страницы
    function initPage(brandData) {
        renderProducts(brandData.products);
        products = brandData.products;
    }

    document.getElementById('reviewForm').addEventListener('submit', function(e) {
        e.preventDefault();

        const text = document.getElementById('review-text').value;

        // Получаем рейтинг
        let rating = 0;
        document.querySelectorAll('#userRating .star').forEach(star => {
            if (star.classList.contains('active')) rating++;
        });

        if (rating === 0) {
            alert('Пожалуйста, поставьте оценку');
            return;
        }

        fetch(`/api/v1/create-review/${productId}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({rating: rating, comment: text})
        }).then(r => r.json()).then(r => {
            if (r.error) {
                alert(r.error);
            } else {
                alert('Спасибо за ваш отзыв!');
                window.location.reload()
            }
        });
    });
});
