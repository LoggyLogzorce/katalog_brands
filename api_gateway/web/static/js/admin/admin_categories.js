document.addEventListener('DOMContentLoaded', () => {
    // DOM
    const tableCtn      = document.getElementById('categoriesTableContainer');
    const addForm       = document.getElementById('addCategoryForm');
    const editModal     = document.getElementById('editCategoryModal');
    const closeModalBtn = editModal.querySelector('.close-modal');
    const cancelBtn     = document.getElementById('cancelEdit');
    const saveBtn       = document.getElementById('saveCategory');
    const searchInput   = document.getElementById('searchInput');

    const addImgInput   = document.getElementById('categoryImage');
    const addImgPrev    = document.getElementById('imagePreview');
    const editImgInput  = document.getElementById('editCategoryImage');
    const editImgPrev   = document.getElementById('editImagePreview');

    let categories = [];
    let currentId  = null;

    // API endpoints
    const API = {
        list:    '/api/v1/categories',
        create:  '/api/v1/admin/category/create',
        update:  id => `/api/v1/admin/category/update/${id}`,
        remove:  id => `/api/v1/admin/category/delete/${id}`
    };

    // Fetch all
    async function fetchCategories() {
        try {
            const res = await fetch(API.list);
            categories = await res.json();
        } catch (e) {
            console.error(e);
            categories = [];
        }
    }

    // Render table
    function render(list = categories) {
        if (!list.length) {
            tableCtn.innerHTML = `<div class="error-message"><i class="fas fa-exclamation-circle"></i><p>Категории не найдены</p></div>`;
            return;
        }
        let html = `<table><thead><tr><th>ID</th><th>Название</th><th>Изображение</th><th>Товаров</th><th>Действия</th></tr></thead><tbody>`;
        list.forEach(c => {
            html += `<tr>
        <td>#${c.category_id}</td>
        <td>${c.name}</td>
        <td><img src="/static/${c.photo}" class="category-image"></td>
        <td>${c.product_count}</td>
        <td>
          <button class="action-btn edit" data-id="${c.category_id}"><i class="fas fa-edit"></i></button>
          <button class="action-btn delete" data-id="${c.category_id}"><i class="fas fa-trash"></i></button>
        </td>
      </tr>`;
        });
        html += `</tbody></table>`;
        tableCtn.innerHTML = html;
        bindTableButtons();
    }

    // Bind row buttons
    function bindTableButtons() {
        document.querySelectorAll('.action-btn.edit').forEach(btn => {
            btn.onclick = () => openEdit(+btn.dataset.id);
        });
        document.querySelectorAll('.action-btn.delete').forEach(btn => {
            btn.onclick = async () => {
                const id = +btn.dataset.id;
                if (!confirm('Удалить категорию?')) return;
                const res = await fetch(API.remove(id), { method: 'DELETE' });
                if (res.ok) {
                    await reload();
                } else {
                    alert('Ошибка при удалении');
                }
            };
        });
    }

    // Open edit modal
    function openEdit(id) {
        currentId = id;
        const c = categories.find(x => x.category_id === id);
        if (!c) return;
        document.getElementById('editCategoryId').value   = c.category_id;
        document.getElementById('editCategoryName').value = c.name;
        editImgPrev.innerHTML = `
      <div class="image-preview">
        <img src="/static/${c.photo}">
      </div>`;
        editModal.style.display = 'flex';
    }

    // Close modal
    function closeModal() {
        editModal.style.display = 'none';
        currentId = null;
    }

    // Preview image helper
    function previewImage(evt, container) {
        const file = evt.target.files[0];
        if (!file) return;
        const reader = new FileReader();
        reader.onload = e => {
            container.innerHTML = `
        <div class="image-preview">
          <img src="${e.target.result}">
        </div>`;
        };
        reader.readAsDataURL(file);
    }

    // Init file-input labels
    document.querySelectorAll('.file-input-container').forEach(c => {
        const label = c.querySelector('.file-input-label');
        const input = c.querySelector('input[type="file"]');
        label.onclick = () => input.click();
    });

    // Search
    searchInput.oninput = () => {
        const term = searchInput.value.trim().toLowerCase();
        render(term
            ? categories.filter(c => c.name.toLowerCase().includes(term))
            : categories);
    };

    // Add category
    addForm.onsubmit = async e => {
        e.preventDefault();
        const name = document.getElementById('categoryName').value.trim();
        if (!name) return alert('Введите название');
        const fd = new FormData(addForm);
        fd.append('category_name', name);
        if (addImgInput.files[0]) fd.append('logoImage', addImgInput.files[0]);
        const res = await fetch(API.create, { method: 'POST', body: fd });
        if (res.ok) {
            addForm.reset();
            addImgPrev.innerHTML = '';
            await reload();
        } else {
            alert('Ошибка при добавлении');
        }
    };

    // Save edit
    saveBtn.onclick = async () => {
        const name = document.getElementById('editCategoryName').value.trim();
        if (!name) return alert('Введите название');
        const fd = new FormData();
        fd.append('category_name', name);
        if (editImgInput.files[0]) fd.append('logoImage', editImgInput.files[0]);
        const res = await fetch(API.update(currentId), { method: 'PUT', body: fd });
        if (res.ok) {
            closeModal();
            await reload();
        } else {
            alert('Ошибка при сохранении');
        }
    };

    // Preview hooks
    addImgInput.onchange  = e => previewImage(e, addImgPrev);
    editImgInput.onchange = e => previewImage(e, editImgPrev);

    // Modal close hooks
    closeModalBtn.onclick = cancelBtn.onclick = closeModal;
    window.onclick = e => { if (e.target === editModal) closeModal(); };

    // Reload data + render
    async function reload() {
        await fetchCategories();
        render();
    }

    // Initial
    (async () => {
        await reload();
    })();
});
