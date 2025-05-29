fetch('/api/brands', {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' }
    })
    .then(res => res.json())
    .then(brands => {
        const container = document.getElementById('brands-list');
        console.log(brands)
        console.log(brands.length)
        if (brands.length === 0) {
            container.innerHTML = '<p style="padding:20px">Брендов не найдено.</p>';
            return;
        }
        brands.forEach(b => {
            const card = document.createElement('div');
            card.className = 'brand-card';
            card.innerHTML = `
            <a href="/brands/${b.Name}">
              <img src="${b.logoURL}" alt="${b.Name}" class="brand-logo">
              <h3>${b.Name}</h3>
            </a>
            <div class="rating">★ ${b.AverageRating.toFixed(1)}</div>
            <p>${b.description}</p>
            <a href="/brands/${b.Name}" class="btn">Подробнее</a>
          `;
            container.appendChild(card);
        });
    })
    .catch(() => {
        document.getElementById('brands-list')
            .innerHTML = '<p style="padding:20px">Ошибка загрузки брендов.</p>';
    });