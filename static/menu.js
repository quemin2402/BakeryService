document.addEventListener('DOMContentLoaded', function () {
  const productContainer = document.querySelector('#products');
  const addProductButton = document.querySelector('#addProductBtn');
  const searchInput = document.querySelector('#searchInput');
  const categorySelect = document.querySelector('#categorySelect');
  const priceMinInput = document.querySelector('#priceMin');
  const priceMaxInput = document.querySelector('#priceMax');
  const filterButton = document.querySelector('#filterButton');
  const searchByIdButton = document.querySelector('#searchButton');
  const searchByIdInput = document.querySelector('#searchProductInput');

  let currentPage = 1;
  const limit = 8; // Количество продуктов на странице

  async function fetchProducts(queryParams = '') {
    try {
      const response = await fetch(`/api/products?page=${currentPage}&limit=${limit}${queryParams}`);
      if (!response.ok) throw new Error('Failed to fetch products');
      const products = await response.json();
      renderProducts(products);

      // Обновить состояние кнопок пагинации
      document.getElementById('currentPage').textContent = `Page ${currentPage}`;
      document.getElementById('prevPage').disabled = currentPage === 1;
      document.getElementById('nextPage').disabled = products.length < limit; // Если меньше лимита, значит, данных больше нет
    } catch (error) {
      console.error('Error:', error);
      alert('Failed to load products. Please try again.');
    }
  }

  document.getElementById('prevPage').addEventListener('click', () => {
    if (currentPage > 1) {
      currentPage--;
      fetchProducts();
    }
  });

  document.getElementById('nextPage').addEventListener('click', () => {
    currentPage++;
    fetchProducts();
  });



  function renderProducts(products) {
    productContainer.innerHTML = '';
    if (products.length === 0) {
      productContainer.innerHTML = '<p class="text-center">No products found.</p>';
      return;
    }
    products.forEach(product => {
      const imagePath = product.image || '/static/images/default.jpg';
      const productCard = `
      <div class="col" data-product-id="${product.id}">
        <div class="card h-100">
          <img src="${imagePath}" class="card-img-top menu-card-img" alt="${product.name}">
          <div class="card-body text-center">
            <h5 class="card-title">${product.name}</h5>
            <p class="card-text price">$${product.price.toFixed(2)}</p>
            <p class="card-text category">${product.category || 'Uncategorized'}</p>
            <button class="btn btn-danger delete-product" data-id="${product.id}">Delete</button>
            <button class="btn btn-primary edit-product" data-id="${product.id}">Edit</button>
          </div>
        </div>
      </div>
    `;
      productContainer.insertAdjacentHTML('beforeend', productCard);
    });
    attachDeleteEvents();
    attachEditEvents();
  }


  addProductButton.addEventListener('click', function () {
    const name = prompt('Enter product name (required):');
    const price = prompt('Enter product price (required):');
    const category = prompt('Enter product category (optional):');
    const description = prompt('Enter description (optional):');
    if (!name || isNaN(price) || price <= 0) {
      alert('Invalid input. Name and Price are required.');
      return;
    }
    createProduct({ name, price: parseFloat(price), category, description });
  });

  async function createProduct(product) {
    try {
      const response = await fetch('/api/product/create', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(product),
      });
      if (!response.ok) throw new Error('Failed to create product');
      alert('Product created successfully!');
      fetchProducts();
    } catch (error) {
      console.error('Error:', error);
      alert('Failed to create product.');
    }
  }

  function attachDeleteEvents() {
    document.querySelectorAll('.delete-product').forEach(button => {
      button.addEventListener('click', async function () {
        const productId = this.getAttribute('data-id');
        if (confirm('Are you sure you want to delete this product?')) {
          await deleteProduct(productId);
        }
      });
    });
  }

  function attachEditEvents() {
    document.querySelectorAll('.edit-product').forEach(button => {
      button.addEventListener('click', async function () {
        const productId = this.getAttribute('data-id');
        const newName = prompt('Enter new product name:');
        const newPrice = prompt('Enter new product price:');
        if (newName && !isNaN(newPrice)) {
          await updateProduct(productId, { name: newName, price: parseFloat(newPrice) });
        } else {
          alert('Invalid input!');
        }
      });
    });
  }

  async function deleteProduct(productId) {
    try {
      const response = await fetch(`/api/product/delete/${productId}`, { method: 'DELETE' });
      if (!response.ok) throw new Error('Failed to delete product');
      alert('Product deleted successfully!');
      fetchProducts();
    } catch (error) {
      console.error('Error:', error);
      alert('Failed to delete product.');
    }
  }

  async function updateProduct(productId, updatedData) {
    try {
      const response = await fetch(`/api/product/update?id=${productId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updatedData),
      });
      if (!response.ok) throw new Error('Failed to update product');
      alert('Product updated successfully!');
      fetchProducts();
    } catch (error) {
      console.error('Error:', error);
      alert('Failed to update product.');
    }
  }

  searchByIdButton.addEventListener('click', async function () {
    const id = searchByIdInput.value.trim();
    if (!id || isNaN(id)) {
      alert('Please enter a valid product ID.');
      return;
    }

    try {
      const response = await fetch(`/api/product/${id}`);
      if (!response.ok) throw new Error('Product not found');
      const product = await response.json();
      renderProducts([product]); // Отобразить только найденный продукт
    } catch (error) {
      console.error('Error:', error);
      alert('Product not found!');
    }
  });


  filterButton.addEventListener('click', function () {
    const searchTerm = searchInput.value.trim();
    const category = categorySelect.value.trim();
    const priceMin = priceMinInput.value.trim() ? parseFloat(priceMinInput.value) : 0;
    const priceMax = priceMaxInput.value.trim() ? parseFloat(priceMaxInput.value) : "";
    const sortBy = document.querySelector('#sortBy').value;
    const sortOrder = document.querySelector('#sortOrder').value;

    let queryParams = `?page=${currentPage}&limit=${limit}`;
    if (searchTerm) queryParams += `&name=${encodeURIComponent(searchTerm)}`;
    if (category) queryParams += `&category=${encodeURIComponent(category)}`;
    if (!isNaN(priceMin)) queryParams += `&priceMin=${priceMin}`;
    if (!isNaN(priceMax)) queryParams += `&priceMax=${priceMax}`;
    if (sortBy) queryParams += `&sortBy=${encodeURIComponent(sortBy)}`;
    if (sortOrder) queryParams += `&sortOrder=${encodeURIComponent(sortOrder)}`;

    // Убираем лишний "&" в конце
    queryParams = queryParams.endsWith('&') ? queryParams.slice(0, -1) : queryParams;

    fetchProducts(queryParams);
  });



  fetchProducts();
  console.log(queryParams); // Убедитесь, что строка запроса правильная

});
