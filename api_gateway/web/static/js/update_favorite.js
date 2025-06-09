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
                if (!res.ok) throw new Error();
                icon.classList.toggle('fas', !isFav);
                icon.classList.toggle('far', isFav);
            } catch {
                const error = new Error('Не удалось обновить избранное');
                console.error(error);
            }
        };
    });
}