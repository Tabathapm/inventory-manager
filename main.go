package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "inventory-manager/database"
    "strconv"
    "strings"
)


func main() {
    // Conexión con la base de datos
    database.InitDB()

    // ---------- PRODUCTOS ---------
    http.HandleFunc("/productos", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        // GET: Pedir lista a la base de datos
        if r.Method == "GET" {
            productos, err := database.ObtenerProductos()
            if err != nil {
                http.Error(w, "Error al leer de la base de datos", http.StatusInternalServerError)
                return
            }
            json.NewEncoder(w).Encode(productos)
            return
        }

        // POST: Guardar en la base de datos
        if r.Method == "POST" {
            var nuevoProducto database.Producto
            
            // Decodificar el JSON que llega
            if err := json.NewDecoder(r.Body).Decode(&nuevoProducto); err != nil {
                http.Error(w, "JSON inválido", http.StatusBadRequest)
                return
            }

            // Guardar en Postgres
            id, err := database.InsertarProducto(nuevoProducto)
            if err != nil {
                http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
                log.Println("Error SQL:", err) // Log para nosotros
                return
            }

            // Asignar el ID nuevo al objeto para devolverlo completo
            nuevoProducto.ID = id
            
            w.WriteHeader(http.StatusCreated)
            json.NewEncoder(w).Encode(nuevoProducto)
        }
    })

    // --------- VENDER ---------
    http.HandleFunc("/vender", func(w http.ResponseWriter, r *http.Request){
        w.Header().Set("Content-Type", "application/json")

        if r.Method == "POST" {
            var nuevaOrden database.OrdenDeCompra

            if err := json.NewDecoder(r.Body).Decode(&nuevaOrden); err != nil {
                http.Error(w, "JSON inválido", http.StatusBadRequest)
                return
            }

            orden, err := database.RegistrarOrdenDeCompra(nuevaOrden)
            if err != nil {
                // Si el error dice "insuficiente" o "no encontrado", es culpa del cliente (400)
                if err.Error() == "stock insuficiente" || err.Error() == "producto no encontrado" {
                    http.Error(w, err.Error(), http.StatusBadRequest) // 400 Bad Request
                } else {
                    // Si es otro error (ej: falló SQL), es culpa nuestra (500)
                    http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
                    log.Println("Error SQL:", err)
                }
                return
            }

            nuevaOrden = orden

            w.WriteHeader(http.StatusCreated)
            json.NewEncoder(w).Encode(nuevaOrden)
        }
    })

    // ----------- UN PRODUCTO X -------------
    http.HandleFunc("/productos/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        if r.Method == "GET" {
            // Si la URL es "/productos/5", esto deja solo el string "5"
            idStr := strings.TrimPrefix(r.URL.Path, "/productos/")

            // Convertir texto a número (int)
            id, err := strconv.Atoi(idStr)
            if err != nil {
                http.Error(w, "El ID debe ser un número válido", http.StatusBadRequest)
                return
            }

            producto, err := database.TraerProducto(id)
            if err != nil {
                if err.Error() == "Producto inexistente" {
                    http.Error(w, "Producto no encontrado", http.StatusNotFound)
                } else {
                    http.Error(w, "Error interno", http.StatusInternalServerError)
                    log.Println("Error SQL:", err)
                }
                return
            }

            json.NewEncoder(w).Encode(producto)
        }
    })

    fmt.Println("Servidor corriendo en el puerto 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}