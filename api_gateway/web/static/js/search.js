document.addEventListener('DOMContentLoaded', function () {
    // Поиск по категориям
    const searchInput = document.querySelector('.search-bar input');
    const searchButton = document.querySelector('.search-bar button');

    searchButton.addEventListener('click', function () {
        if (searchInput.value.trim() !== '') {
            window.location.href = `/search?query=${encodeURIComponent(searchInput.value)}`;
        } else {
            searchInput.focus();
        }
    });

    searchInput.addEventListener('keypress', function (e) {
        if (e.key === 'Enter') searchButton.click();
    });
})