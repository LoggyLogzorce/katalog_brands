document.addEventListener('DOMContentLoaded', () => {
    const endpoints = {
        products: '/api/v1/admin/products',
        update_product: '/api/v1/admin/product/update',
        create_product: '/api/v1/admin/product/create',
        delete_product: '/api/v1/admin/product/delete',
        brands: '/api/v1/admin/brands',
        categories: '/api/v1/categories'
    };

    let products = [];
    let filtered = [];

    // DOM элементы
    const tbody = document.querySelector('table tbody');
    const statusFilter = document.getElementById('statusFilter');
    const brandFilter = document.getElementById('brandFilter');
    const categoryFilter = document.getElementById('categoryFilter');
    const sortFilter = document.getElementById('sortFilter');
    const searchInput = document.querySelector('.search-container input');
    const addProductBtn = document.getElementById('addProductBtn');
    const applyFiltersBtn = document.getElementById('applyFilters');

    // Модалка и форма
    const editModal = document.getElementById('editProductModal');
    const modalTitle = editModal.querySelector('.modal-title');
    const closeModalBtn = editModal.querySelector('.close-modal');
    const cancelBtn = document.getElementById('cancelEdit');
    const saveBtn = document.getElementById('saveProduct');
    const form = document.getElementById('productEditForm');

    // Поля формы
    const fName = document.getElementById('productName');
    const fPrice = document.getElementById('productPrice');
    const fBrand = document.getElementById('productBrand');
    const fStatus = document.getElementById('productStatus');
    const fDesc = document.getElementById('productDescription');
    const previewContainer = document.getElementById('previewContainer');
    const fileInput = document.getElementById('productImages');
    const fCategory = document.getElementById('productCategory');

    Promise.all([
        fetch(endpoints.products).then(r => r.json()),
        fetch(endpoints.brands).then(r => r.json()),
        fetch(endpoints.categories).then(r => r.json())
    ]).then(([prd, brd, cat]) => {
        products = prd;
        initFilters(brd, cat);
        initModalSelects(brd, cat);
        applyFilters();
        bindActions();
    }).catch(console.error);

    function initFilters(brands, categories) {
        brands.forEach(b => brandFilter.insertAdjacentHTML('beforeend', `<option value="${b.id}">${b.name}</option>`));
        categories.forEach(c => categoryFilter.insertAdjacentHTML('beforeend', `<option value="${c.category_id}">${c.name}</option>`));
    }

    function initModalSelects(brands, categories) {
        // бренды
        fBrand.innerHTML = '';
        brands.forEach(b => fBrand.insertAdjacentHTML('beforeend', `<option value="${b.id}">${b.name}</option>`));
        // категории
        fCategory.innerHTML = '';
        categories.forEach(c => fCategory.insertAdjacentHTML('beforeend', `<option value="${c.category_id}">${c.name}</option>`));
    }

    function bindActions() {
        applyFiltersBtn.onclick = applyFilters;
        searchInput.addEventListener('input', applyFilters);
        closeModalBtn.addEventListener('click', closeModal);
        cancelBtn.addEventListener('click', closeModal);
        window.addEventListener('click', e => e.target === editModal && closeModal());
        saveBtn.addEventListener('click', saveProduct);
        fileInput.addEventListener('change', handleFiles);
        addProductBtn.addEventListener('click', () => {
            modalTitle.textContent = 'Добавить товар';
            form.reset();
            previewContainer.innerHTML = '';
            editModal.style.display = 'flex';
        });
        window.addEventListener('keydown', e => {
            if (e.key === 'Escape' && editModal.style.display === 'flex') {
                closeModal();
            }
        });
    }


    function applyFilters() {
        const status = statusFilter.value;
        const brandId = brandFilter.value;
        const catId = categoryFilter.value;
        const sort = sortFilter.value;
        const query = searchInput.value.trim().toLowerCase();

        filtered = products
            .filter(p => !status || p.status === status)
            .filter(p => !brandId || String(p.brand_id) === brandId)
            .filter(p => !catId || String(p.category_id) === catId)
            .filter(p => !query || String(p.product_id).includes(query) || p.name.toLowerCase().includes(query));

        filtered.sort((a, b) => {
            switch (sort) {
                case 'newest': return new Date(b.created_at) - new Date(a.created_at);
                case 'oldest': return new Date(a.created_at) - new Date(b.created_at);
                case 'price-low': return a.price - b.price;
                case 'price-high': return b.price - a.price;
                default: return 0;
            }
        });

        renderTable();
    }

    function renderTable() {
        tbody.innerHTML = '';
        if (!filtered.length) {
            tbody.innerHTML = '<tr><td colspan="7" style="text-align:center">Ничего не найдено</td></tr>';
            return;
        }
        filtered.forEach(p => {
            const date = new Date(p.created_at).toLocaleDateString('ru-RU');
            let cls = {'pending': 'status-pending', 'approved': 'status-approved', 'rejected': 'status-rejected'}[p.status];
            let txt = {'pending': 'На модерации', 'approved': 'Одобрено', 'rejected': 'Отклонено'}[p.status];
            tbody.insertAdjacentHTML('beforeend', `
                <tr>
                    <td>#${String(p.product_id).padStart(3, '0')}</td>
                    <td>${p.name}</td>
                    <td><span class="price-tag">${p.price.toLocaleString('ru-RU')} ₽</span></td>
                    <td>${p.brand.name}</td>
                    <td>${date}</td>
                    <td><span class="status ${cls}">${txt}</span></td>
                    <td>
                        <button class="action-btn edit" data-id="${p.product_id}" title="Редактировать"><i class="fas fa-edit"></i></button>
                        <button class="action-btn delete" data-id="${p.product_id}" title="Удалить"><i class="fas fa-trash"></i></button>
                    </td>
                </tr>`);
        });
        bindRowActions();
    }

    function bindRowActions() {
        document.querySelectorAll('.action-btn.edit').forEach(btn => btn.onclick = () => openEditModal(+btn.dataset.id));
        document.querySelectorAll('.action-btn.delete').forEach(btn => btn.onclick = () => deleteProduct(+btn.dataset.id));
    }

    function openEditModal(id) {
        const p = products.find(x => x.product_id === id);
        if (!p) return;

        modalTitle.textContent = `Редактирование товара #${id}`;
        fName.value = p.name;
        fPrice.value = p.price;
        fBrand.value = p.brand_id;
        fCategory.value = p.category_id;
        fStatus.value = p.status;
        fDesc.value = p.description;

        previewContainer.innerHTML = '';
        p.product_urls.forEach(u => {
            const div = document.createElement('div');
            div.className = 'image-preview';
            div.innerHTML = `<img src="/static/${u.url}" alt="">`;
            previewContainer.append(div);
        });

        editModal.style.display = 'flex';
    }

    function closeModal() {
        editModal.style.display = 'none';
    }

    function saveProduct() {
        const match = modalTitle.textContent.match(/#(\d+)/);
        const isEdit = Boolean(match);
        const id = isEdit ? +match[1] : null;
        if (!form.checkValidity()) return form.reportValidity();

        const url = isEdit ? `${endpoints.update_product}/${id}` : endpoints.create_product;
        const method = isEdit ? 'PUT' : 'POST';
        const fd = new FormData(form);

        fetch(url, {
            method,
            body: fd
        })
            .then(r => {
                if (r.ok) window.location.reload();
            })
            .catch(console.error);
    }

    function deleteProduct(id) {
        if (!confirm(`Удалить товар #${id}?`)) return;
        fetch(`${endpoints.delete_product}/${id}`, { method: 'DELETE' })
            .then(r => {
                if (r.ok) {
                    window.location.reload()
                }
            });
    }

    function handleFiles() {
        previewContainer.innerHTML = '';
        Array.from(fileInput.files).forEach(file => {
            const reader = new FileReader();
            reader.onload = e => {
                const div = document.createElement('div'); div.className = 'image-preview';
                div.innerHTML = `<img src="${e.target.result}" alt="">`;
                previewContainer.append(div);
            };
            reader.readAsDataURL(file);
        });
    }
});
