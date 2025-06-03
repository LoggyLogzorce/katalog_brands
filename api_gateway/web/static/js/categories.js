fetch('/api/v1/categories', {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' }
})
    .then(res => res.json())
    .then(categories => {
        const container = document.getElementById('categories-list');
        if (categories.length === 0) {
            container.innerHTML = '<p style="padding:20px">Категорий не найдено.</p>';
            return;
        }
        categories.forEach(cat => {
            const card = document.createElement('div');
            card.className = 'category-card';
            card.innerHTML = `
            <a href="/categories/${cat.id}">
              <h3>${cat.name}</h3>
            </a>
            <p>${cat.description}</p>
            <a href="/categories/${cat.id}" class="btn">Перейти</a>
          `;
            container.appendChild(card);
        });
    })
    .catch(() => {
        document.getElementById('categories-list')
            .innerHTML = '<p style="padding:20px">Ошибка загрузки категорий.</p>';
    });