package database // El nombre del paquete debe coincidir con la carpeta

import (
    "database/sql"
    "fmt"
    "log"
	"time"

    _ "github.com/lib/pq" // Importar el driver acá porque acá se usa
)

type Producto struct {
    ID     int     `json:"id"`
    Nombre string  `json:"nombre"`
    Precio float64 `json:"precio"`
    Stock  int     `json:"stock"`
}

type OrdenDeCompra struct {
	ID          int `json:"id"`
	Producto_id int `json:"Producto_id"`
	Cantidad    int `json:"cantidad"`
	Total      float64   `json:"total"`
	Fecha      time.Time `json:"fecha"`
}

// DB es pública (Mayúscula) para que en main.go se pueda usar
var DB *sql.DB

// InitDB inicia la conexión (Mayúscula para exportarla)
func InitDB() {
    connStr := "user=postgres password=secret dbname=inventory_db sslmode=disable host=db"

    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatal("No se pudo conectar a la BDD:", err)
    }
    
    fmt.Println("🚀 ¡Conexión a PostgreSQL exitosa (desde el paquete database)!")

    crearTablas()
}

// crearTablas es privada (minúscula) porque solo se llama desde InitDB, 
// a main.go no le interesa saber cómo se crean las tablas.
func crearTablas() {
    queryProductos := `
    CREATE TABLE IF NOT EXISTS productos (
        id SERIAL PRIMARY KEY,
        nombre TEXT NOT NULL,
        precio NUMERIC(10, 2) NOT NULL DEFAULT 0,
        stock INTEGER NOT NULL DEFAULT 0
    );`

    if _, err := DB.Exec(queryProductos); err != nil {
        log.Fatal("Error creando tabla productos:", err)
    }

	fmt.Println("✅ Tabla 'productos' lista.")

	// ---------------------------------------
	queryOrdenDeCompra := `
	CREATE TABLE IF NOT EXISTS ordenes_compra (
		id SERIAL PRIMARY KEY,
		producto_id INTEGER REFERENCES productos(id), 
		cantidad INTEGER NOT NULL DEFAULT 1,
		total NUMERIC(10, 2) NOT NULL DEFAULT 0,
       	fecha TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := DB.Exec(queryOrdenDeCompra); err != nil {
		log.Fatal("Error creando tabla ordenes_compra:", err)
	}

	fmt.Println("✅ Tabla 'ordenes_compra' lista.")
}

// --- NUEVAS FUNCIONES SQL ---

// Función para INSERTAR un producto en la base de datos
func InsertarProducto(producto Producto) (int, error) {
    // Usar $1, $2, $3 para evitar inyecciones SQL (seguridad básica)
    query := `INSERT INTO productos (nombre, precio, stock) VALUES ($1, $2, $3) RETURNING id`
    
    var id int
    // Ejecutar la query y nos devuelve el ID generado
    err := DB.QueryRow(query, producto.Nombre, producto.Precio, producto.Stock).Scan(&id)
    if err != nil {
        return 0, err
    }
    return id, nil
}

// Función para LEER todos los productos
func ObtenerProductos() ([]Producto, error) {
    rows, err := DB.Query("SELECT id, nombre, precio, stock FROM productos")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var productos []Producto

    // Recorremos fila por fila lo que nos devolvió Postgres
    for rows.Next() {
        var producto Producto
        // "Escaneamos" los datos de la fila a nuestra variable producto
        if err := rows.Scan(&producto.ID, &producto.Nombre, &producto.Precio, &producto.Stock); err != nil {
            return nil, err
        }
        productos = append(productos, producto)
    }
    return productos, nil
}

// Función para buscar un determinado producto
func TraerProducto(id int) (Producto, error) {
    // Consulta
    query := "SELECT id, nombre, precio, stock FROM productos WHERE id = $1"
    
    var producto Producto

    // QueryRow() y .Scan() 
	// QueryRow está diseñada específicamente para buscar uno solo
    err := DB.QueryRow(query, id).Scan(&producto.ID, &producto.Nombre, &producto.Precio, &producto.Stock)

    if err != nil {
        // Si no existe, se devuelve el producto vacío y el error
        return Producto{}, fmt.Errorf("Producto inexistente")
    }

    return producto, nil
}

// Función para restar stock
func RestarStock(id int, cantidad int) (Producto, error) {
    
    producto, err := TraerProducto(id)
    if err != nil {
        return Producto{}, err // Si no existe o falla la BD, salimos aquí
    }

    // Verificación de stock
    if producto.Stock < cantidad {
        // Error nuevo personalizado
        return Producto{}, fmt.Errorf("stock insuficiente: hay %d, pides %d", producto.Stock, cantidad)
    }

    // Consulta
    query := "UPDATE productos SET stock = stock - $1 WHERE id = $2"

	// Query: Es para cuando se espera que la base de datos responda con filas (SELECT).
	// Exec: Es para cuando solo se da una orden (UPDATE, DELETE, INSERT) y no se esperan filas de vuelta. 

    // Se pasa 'cantidad' primero porque es el $1
    _, err = DB.Exec(query, cantidad, id)
    if err != nil {
        return Producto{}, err
    }

    // Actualizar la variable local para devolver el dato real
    producto.Stock = producto.Stock - cantidad

    return producto, nil
}

// Registrar orden de compra
func RegistrarOrdenDeCompra(orden OrdenDeCompra) (OrdenDeCompra, error) {
    
	producto, err := TraerProducto(orden.Producto_id)
	if err != nil {
        return OrdenDeCompra{}, fmt.Errorf("producto no encontrado")
    }

	if producto.Stock < orden.Cantidad {
		return OrdenDeCompra{}, fmt.Errorf("stock insuficiente")
	}

	orden.Total = producto.Precio * float64(orden.Cantidad)
	_, err = RestarStock(orden.Producto_id, orden.Cantidad)
	if err != nil {
		return OrdenDeCompra{}, fmt.Errorf("error al actualizar stock")
	}
	
	query := `INSERT INTO ordenes_compra (producto_id, cantidad, total) VALUES ($1, $2, $3) RETURNING id, fecha`

    // Ejecutamos y recuperamos el ID generado y la FECHA automática
	err = DB.QueryRow(query, orden.Producto_id, orden.Cantidad, orden.Total).Scan(&orden.ID, &orden.Fecha)
	if err != nil {
		return OrdenDeCompra{}, err
	}

	return orden, nil
}

