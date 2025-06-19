let brands = [];
let currentBrandId = null;

const brandsTableContainer = document.getElementById('brandsTableContainer');
const editModal = document.getElementById('editBrandModal');
const closeModalBtn = editModal.querySelector('.close-modal');
const cancelBtn = document.getElementById('cancelEdit');
const saveBrandBtn = document.getElementById('saveBrand');
const addBrandBtn = document.getElementById('addBrandBtn');
const applyFiltersBtn = document.getElementById('applyFilters');
const searchInput = document.getElementById('searchInput');

async function fetchBrands() {
    try {
        const res = await fetch('/api/v1/admin/brands');
        return await res.json();
    } catch (e) {
        console.error(e);
        return [];
    }
}

function formatDate(iso) {
    return new Date(iso).toLocaleDateString('ru-RU');
}

function renderBrands(list) {
    if (!list.length) {
        brandsTableContainer.innerHTML = `<div class="error-message"><i class="fas fa-exclamation-circle"></i><p>Бренды не найдены</p></div>`;
        return;
    }
    let html = `<table>
        <thead><tr>
            <th>ID</th><th>Название</th><th>Создатель</th><th>Дата</th><th>Статус</th><th>Действия</th>
        </tr></thead><tbody>`;
    list.forEach(b => {
        let cls = {'pending': 'status-pending', 'approved': 'status-approved', 'rejected': 'status-rejected'}[b.status];
        let txt = {'pending': 'На модерации', 'approved': 'Одобрено', 'rejected': 'Отклонено'}[b.status];
        html += `<tr>
            <td>#${b.id}</td>
            <td>${b.name}</td>
            <td>${b.creator_id}</td>
            <td>${formatDate(b.created_at)}</td>
            <td><span class="status ${cls}">${txt}</span></td>
            <td>
                <button class="action-btn edit" data-id="${b.id}"><i class="fas fa-edit"></i></button>
                <button class="action-btn delete" data-id="${b.id}"><i class="fas fa-trash"></i></button>
            </td>
        </tr>`;
    });
    html += `</tbody></table>`;
    brandsTableContainer.innerHTML = html;
    addListeners();
}

function addListeners() {
    document.querySelectorAll('.action-btn.edit').forEach(btn => {
        btn.onclick = () => openEditModal(+btn.dataset.id);
    });
    document.querySelectorAll('.action-btn.delete').forEach(btn => {
        btn.onclick = () => {
            const id = +btn.dataset.id;
            if (confirm('Удалить бренд?')) {
                fetch(`/api/v1/admin/brand/${id}`, {method: 'DELETE'})
                    .then(r => {
                        if (r.ok) {
                            brands = brands.filter(b => b.id !== id);
                            renderBrands(brands);
                        }
                    });
            }
        };
    });
}

function openEditModal(id) {
    currentBrandId = id;
    const b = brands.find(x => x.id === id);
    document.getElementById('brandId').value = b.id;
    document.getElementById('brandName').value = b.name;
    document.getElementById('brandOwner').value = b.creator_id;
    document.getElementById('brandStatus').value = b.status;
    document.getElementById('brandDescription').value = b.description;
    document.getElementById('logoPreview').innerHTML = `
        <div class="image-preview">
            <img src="/static/${b.logo_url}" alt="">
        </div>`;
    editModal.style.display = 'flex';
}

function closeModal() {
    editModal.style.display = 'none';
}

saveBrandBtn.onclick = () => {
    const isEdit = Boolean(currentBrandId);
    const url = isEdit
        ? `/api/v1/admin/brand/${encodeURIComponent(currentBrandId)}`
        : `/api/v1/admin/brand/create`;
    const method = isEdit ? 'PUT' : 'POST';
    const form = document.getElementById('brandEditForm');
    if (!form.checkValidity()) return form.reportValidity();
    const fd = new FormData(form);
    fetch(url, {
        method: method,
        body: fd
    }).then(r => {
        if (r.ok) window.location.reload();
    });
};

addBrandBtn.onclick = () => {
    currentBrandId = null;
    document.querySelector('.modal-title').textContent = 'Добавить бренд';
    document.getElementById('brandEditForm').reset();
    document.getElementById('logoPreview').innerHTML = '';
    editModal.style.display = 'flex';
};

closeModalBtn.onclick = cancelBtn.onclick = closeModal;
window.onclick = e => {
    if (e.target === editModal) closeModal();
};
window.addEventListener('keydown', e => {
    if (e.key === 'Escape' && editModal.style.display === 'flex') {
        closeModal();
    }
});

applyFiltersBtn.onclick = filterAndSort;
searchInput.oninput = filterAndSort;

function filterAndSort() {
    let filtered = brands.filter(b => {
        const st = document.getElementById('statusFilter').value;
        const term = searchInput.value.toLowerCase();
        return (!st || b.status === st) &&
            (!term || b.name.toLowerCase().includes(term));
    });
    const sort = document.getElementById('sortFilter').value;
    if (sort === 'newest') filtered.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
    if (sort === 'oldest') filtered.sort((a, b) => new Date(a.created_at) - new Date(b.created_at));
    if (sort === 'name') filtered.sort((a, b) => a.name.localeCompare(b.name));
    renderBrands(filtered);
}

document.addEventListener('DOMContentLoaded', async () => {
    brands = await fetchBrands();
    renderBrands(brands);
});

// Превью логотипа
const logoInput = document.getElementById('brandLogo');
const logoPreview = document.getElementById('logoPreview');

logoInput.addEventListener('change', () => {
    logoPreview.innerHTML = '';
    const file = logoInput.files[0];
    if (!file || !file.type.startsWith('image/')) return;
    const reader = new FileReader();
    reader.onload = e => {
        const div = document.createElement('div');
        div.className = 'image-preview';
        div.innerHTML = `
            <img src="${e.target.result}" alt="Логотип">
            <button class="remove-image">&times;</button>
        `;
        logoPreview.append(div);

        // Удаление превью и сброс input
        div.querySelector('.remove-image').onclick = () => {
            logoInput.value = '';
            logoPreview.innerHTML = '';
        };
    };
    reader.readAsDataURL(file);
});

