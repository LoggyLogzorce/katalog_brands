function attachFavoriteHandlers() {
    document.querySelectorAll('.favorite-btn').forEach(btn => {
        btn.onclick = async function (e) {
            e.stopPropagation();
            const productId = btn.dataset.id;
            const icon = btn.querySelector('i');
            const isFav = icon.classList.contains('fas');
            try {
                const method = isFav ? 'DELETE' : 'POST';
                const res = await fetch(`/api/v1/favorites/${productId}`, {
                    method,
                });
                if (!res.ok) {
                    const errJson = await res.json();
                    alert(errJson.error || 'Ошибка при обновлении избранного');
                    throw new Error(errJson.error);
                }
                icon.classList.toggle('fas', !isFav);
                icon.classList.toggle('far', isFav);
            } catch (err) {
                console.error('Не удалось обновить избранное:', err);
            }
        };
    });
}