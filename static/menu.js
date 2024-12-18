document.addEventListener('DOMContentLoaded', function () {
  const productContainer = document.querySelector('#products');

  async function fetchProducts() {
    try {
      const response = await fetch('/api/products');
      if (!response.ok) {
        throw new Error('Failed to fetch products');
      }
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
          </div>
        </div>
      </div>
    `;
      productContainer.insertAdjacentHTML('beforeend', productCard);
    });
    attachDeleteEvents();
  }


  function attachDeleteEvents() {
    document.querySelectorAll('.delete-product').forEach(button => {
      button.addEventListener('click', async function () {
        const productId = this.getAttribute('data-id');
        await deleteProduct(productId);
      });
    });
  }

  async function deleteProduct(productId) {
    try {
      const response = await fetch(`/api/product/delete/${productId}`, {
        method: 'DELETE',
      });

      if (!response.ok) {
        throw new Error('Failed to delete product');
      }

      alert('Product deleted successfully!');
      fetchProducts();
    } catch (error) {
      console.error('Error:', error);
      alert('Failed to delete product. Please try again.');
    }
  }

  fetchProducts();
});
