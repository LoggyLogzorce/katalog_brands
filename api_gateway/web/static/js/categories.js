document.addEventListener('DOMContentLoaded', function() {
    // Поиск по категориям
    const searchInput = document.querySelector('.search-bar input');
    const searchButton = document.querySelector('.search-bar button');

    searchButton.addEventListener('click', function() {
        if (searchInput.value.trim() !== '') {
            alert(`Поиск: "${searchInput.value}"\nРезультаты будут показаны на следующей странице.`);
            // В реальном приложении здесь будет переход на страницу поиска
        } else {
            searchInput.focus();
        }
    });

    searchInput.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            searchButton.click();
        }
    });

    // Анимация при наведении на категории
    const categoryCards = document.querySelectorAll('.category-card');
    categoryCards.forEach(card => {
        card.addEventListener('mouseenter', function() {
            this.querySelector('.category-link').style.gap = '12px';
        });

        card.addEventListener('mouseleave', function() {
            this.querySelector('.category-link').style.gap = '8px';
        });
    });

    // Клик по категории
    categoryCards.forEach(card => {
        card.addEventListener('click', function(e) {
            if (!e.target.closest('.category-link')) {
                const link = this.querySelector('.category-link');
                alert(`Переход в категорию: "${this.querySelector('.category-name').textContent}"`);
                // В реальном приложении: window.location.href = link.href;
            }
        });
    });
});