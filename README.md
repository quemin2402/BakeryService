# Bakery Service

## Description of the project

**Bakery Service** – is a web application for managing bakery processes. The goal of the project is to automate the work with orders, menus and customers in order to improve the user experience and increase the efficiency of the bakery.

The main functionality:
- Receiving and processing orders via a web interface.
- Managing a database of customers, products, and orders.
- Seamless interaction between the **frontend** (HTML/CSS/JS) and **Golang** server using JSON-based RESTful API.
- Dynamic rendering of product data fetched from the server.
- Product management: adding, updating, and deleting items in the menu.

## Team members
- Merey Ibraim (Group SE-2308)

## Screenshots of the website

<img width="1440" alt="The main page" src="https://github.com/user-attachments/assets/af1d8c63-b438-463a-8fb6-c695cd083df9" />

<img width="1435" alt="The menu page" src="https://github.com/user-attachments/assets/9b8e82a6-c69b-4c98-9a08-0185b0358bb4" />


## How to launch a project

Follow the steps below to launch the project:

### **1. Clone the repository**
```bash
git clone https://github.com/yourusername/bakery-service.git
cd bakery-service
```

### **2. Install dependencies**
Ensure you have **Golang** and **PostgreSQL** installed on your computer.

- **Install Go dependencies**:
```bash
go mod tidy
```

### **3. Set up PostgreSQL database**
1. Create a database named `bakery`:
   ```sql
   CREATE DATABASE bakery;
   ```
2. Update `db/connection.go` with your database credentials:
   ```go
   connStr := "host=localhost user=yourusername password=yourpassword dbname=bakery port=5432 sslmode=disable"
   ```
3. Run the project to automatically apply the migrations:
```bash
go run main.go
```

### **4. Run the server**
Start the Golang server:
```bash
go run main.go
```
The server will be available at `http://localhost:8080`.

### **5. Launch the web interface**
1. Open your browser and go to:
   ```
   http://localhost:8080/static/menu.html
   ```

## A tool and a resource
- **Programming language**: Golang  
- **Database**: PostgreSQL  
- **Frontend**: HTML, CSS, JavaScript  
- **API Testing**: Postman  
- **ORM**: GORM  
- **Web Framework**: Gorilla Mux  
- **Additional Tools**: Golang Migrate (for database migrations)

### **Example endpoints for testing the API**:

1. **Get all products**:  
   - Method: `GET`  
   - Endpoint: `/api/products`  

2. **Get product by ID**:  
   - Method: `GET`  
   - Endpoint: `/api/product/{id}`  

3. **Create a new product**:  
   - Method: `POST`  
   - Endpoint: `/api/product/create`  
   - Body (JSON):
     ```json
     {
       "name": "French Baguette",
       "price": 3.5,
       "description": "Classic French bread."
     }
     ```

4. **Update product by ID**:  
   - Method: `PUT`  
   - Endpoint: `/api/product/update?id={id}`  
   - Body (JSON):
     ```json
     {
       "name": "Updated Baguette",
       "price": 3.8,
       "description": "Updated description."
     }
     ```

5. **Delete product by ID**:  
   - Method: `DELETE`  
   - Endpoint: `/api/product/delete/{id}`  
