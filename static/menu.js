document.addEventListener('DOMContentLoaded', function () {
  const productContainer = document.querySelector('#products');
  const addProductButton = document.querySelector('#addProductBtn');
  const searchInput = document.querySelector('#searchInput');
  const searchButton = document.querySelector('#searchButton');

  async function fetchProducts() {
    try {
      const response = await fetch('/api/products');
      if (!response.ok) throw new Error('Failed to fetch products');
      const products = await response.json();
      renderProducts(products);
    } catch (error) {
      console.error('Error:', error);
      alert('Failed to load products. Please try again.');
    }
  }

  function renderProducts(products) {
    productContainer.innerHTML = '';
    products.forEach(product => {
      const imagePath = product.image || '/static/images/default.jpg';
      const productCard = `
        <div class="col" data-product-id="${product.id}">
          <div class="card h-100">
            <img src="${imagePath}" class="card-img-top menu-card-img" alt="${product.name}">
            <div class="card-body text-center">
              <h5 class="card-title">${product.name}</h5>
              <p class="card-text price">$${product.price.toFixed(2)}</p>
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
    const description = prompt('Enter description (optional):');
    if (!name || isNaN(price) || price <= 0) {
      alert('Invalid input. Name and Price are required.');
      return;
    }
    createProduct({ name, price: parseFloat(price), description });
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

  searchButton.addEventListener('click', async function () {
    const id = searchInput.value.trim();
    if (!id || isNaN(id)) {
      alert('Please enter a valid product ID.');
      return;
    }
    try {
      const response = await fetch(`/api/product/${id}`);
      if (!response.ok) throw new Error('Product not found');
      const product = await response.json();
      renderProducts([product]);
    } catch (error) {
      console.error('Error:', error);
      alert('Product not found!');
    }
  });

  fetchProducts();
});
