// Инициализация функционала страницы
document.addEventListener('DOMContentLoaded', function() {
    // Обработчики для кнопок действий
    const actionButtons = document.querySelectorAll('.action-btn');
    actionButtons.forEach(button => {
        button.addEventListener('click', function() {
            if (this.querySelector('.fa-heart')) {
                const icon = this.querySelector('.fa-heart');
                if (icon.classList.contains('far')) {
                    icon.classList.remove('far');
                    icon.classList.add('fas');
                    icon.style.color = '#FFB6C1';
                } else {
                    icon.classList.remove('fas');
                    icon.classList.add('far');
                    icon.style.color = '';
                }
            } else if (this.querySelector('.fa-shopping-cart')) {
                alert('Товар добавлен в корзину!');
            }
        });
    });

    // Обработчики для кнопок профиля
    const profileButtons = document.querySelectorAll('.btn');
    profileButtons.forEach(button => {
        button.addEventListener('click', function(e) {
            if (!this.href) {
                e.preventDefault();
                alert('Функционал в разработке');
            }
        });
    });
});