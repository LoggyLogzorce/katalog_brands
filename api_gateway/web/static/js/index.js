// Создаёт DOM-элемент для карточки бренда
function createBrandCard(item) {
    const div = document.createElement('div');
    div.className = 'brand-card';
    div.innerHTML = `
      <img src="${item.logoURL}" alt="${item.name}" class="brand-logo">
      <h3>${item.name}</h3>
      <div class="rating">${'★'.repeat(Math.floor(item.averageRating))}${'☆'.repeat(5 - Math.floor(item.averageRating))} (${item.averageRating.toFixed(1)})</div>
      <p>${item.description}</p>
      <button class="btn btn-primary" onclick="location.href='/brands/${item.id}'">Подробнее</button>
    `;
    return div;
}

document.addEventListener('DOMContentLoaded', () => {
    // 1. Загрузка популярных брендов
    fetch('/api/v1/brands/popular', {
        method: "GET"
    })
        .then(res => {
            if (!res.ok) throw new Error('Не удалось получить популярные бренды');
            return res.json();
        })
        .then(data => {
            const grid = document.getElementById('popular-brands-grid');
            grid.innerHTML = ''; // убираем статический CSS-скрытый контент
            if (Array.isArray(data) && data.length) {
                data.forEach(item => {
                    const card = createBrandCard(item);
                    grid.appendChild(card);
                });
            } else {
                grid.innerHTML = '<p style="padding:20px">Нет популярных брендов.</p>';
            }
        })
        .catch(err => {
            console.error(err);
            const grid = document.getElementById('popular-brands-grid');
            grid.innerHTML = '<p style="padding:20px">Ошибка загрузки популярных брендов.</p>';
        });

    // 2. Загрузка новых брендов
    fetch('/api/v1/brands/new', {
        method: "GET"
    })
        .then(res => {
            if (!res.ok) throw new Error('Не удалось получить новые бренды');
            return res.json();
        })
        .then(data => {
            const grid = document.getElementById('new-brands-grid');
            grid.innerHTML = ''; // убираем «Загрузка…»
            if (Array.isArray(data) && data.length) {
                data.forEach(item => {
                    const card = createBrandCard(item);
                    grid.appendChild(card);
                });
            } else {
                grid.innerHTML = '<p style="padding:20px">Нет новых брендов.</p>';
            }
        })
        .catch(err => {
            console.error(err);
            const grid = document.getElementById('new-brands-grid');
            grid.innerHTML = '<p style="padding:20px">Ошибка загрузки новых брендов.</p>';
        });
});