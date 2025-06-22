document.addEventListener('DOMContentLoaded', () => {
    const input = document.querySelector('.search-input');
    const btn = document.querySelector('.search-btn');

    function goSearch() {
        const q = input.value.trim();
        if (!q) return;
        // Перенаправляем на страницу результатов поиска
        window.location.href = `/search?q=${encodeURIComponent(q)}`;
    }

    // По клику на кнопку
    btn.addEventListener('click', goSearch);

    // По Enter в поле ввода
    input.addEventListener('keydown', e => {
        if (e.key === 'Enter') {
            e.preventDefault();
            goSearch();
        }
    });
});