package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
)

type Producto struct {
    ID     int     `json:"id"`
    Nombre string  `json:"nombre"`
    Precio float64 `json:"precio"`
    Stock  int     `json:"stock"`
}

type OrdenDeCompra struct {
	ProductoID int `json:"producto_id"`
	Cantidad int `json:"cantidad"`
}

type Respuesta struct {
	Mensaje string `json:"mensaje"`
}

var inventario = []Producto{
    {ID: 1, Nombre: "Mouse Inalámbrico", Precio: 15000.50, Stock: 50},
    {ID: 2, Nombre: "Teclado Mecánico", Precio: 45000.00, Stock: 20},
    {ID: 3, Nombre: "Monitor 24 pulgadas", Precio: 120000.00, Stock: 5},
}

var mu sync.Mutex

func main() {
    http.HandleFunc("/productos", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        respuesta := Respuesta{}

        // OPCIÓN 1: Alguien quiere LEER la lista (GET)
        if r.Method == "GET" {
            json.NewEncoder(w).Encode(inventario)
            return
        }

        // OPCIÓN 2: Alguien quiere AGREGAR un producto (POST)
        if r.Method == "POST" {
            var nuevoProducto Producto
            
            // json.NewDecoder hace lo inverso: 
            // Lee el "Cuerpo" (Body) de la petición 'r' y lo decodifica dentro de la variable 'nuevoProducto'
            err := json.NewDecoder(r.Body).Decode(&nuevoProducto)
            if err != nil {
                respuesta.Mensaje = "Error al leer el JSON"
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(respuesta)
                return
            }

            mu.Lock()
            defer mu.Unlock()

            // Agregamos el producto al slice global (append)
            inventario = append(inventario, nuevoProducto)

            // Respondemos con el producto creado para confirmar (Status 201 Created es lo correcto)
            w.WriteHeader(http.StatusCreated)
            json.NewEncoder(w).Encode(nuevoProducto)
            return
        }
    })

	// VENDER
	http.HandleFunc("/vender", func (w http.ResponseWriter, r *http.Request)  {
		w.Header().Set("Content-Type", "application/json")

		respuesta := Respuesta{}

		// sólo se aceptan métodos POST
		if r.Method != "POST" {
			respuesta.Mensaje = "Método no permitido"

            /*  Al usar "http.Error", Go ignora el struct y lo manda como texto plano
                y el front espera un JSON, no texto plano
			http.Error(w, respuesta.Mensaje, http.StatusMethodNotAllowed)  */
            w.WriteHeader(http.StatusMethodNotAllowed)
            json.NewEncoder(w).Encode(respuesta)
            
			return
		}

		var orden OrdenDeCompra
        // se decodifica lo que nos manda el usuario (qué quiere comprar y cuánto)
        err := json.NewDecoder(r.Body).Decode(&orden)
        if err != nil {
			respuesta.Mensaje = "JSON inválido"

            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(respuesta)
            return
        }
		
		productoEncontrado := false

        // antes de buscar en el inventario, hay que cerrar la puerta para que otro no entre
        mu.Lock()

        // defer mu.Unlock() asegura que la puerta se abra siempre al final de esta función pase lo que pase
        defer mu.Unlock()


		// se busca el producto en el inventario
		for producto := range inventario {
			if inventario[producto].ID == orden.ProductoID {
				productoEncontrado = true
				// si hay suficiente stock
				if inventario[producto].Stock >= orden.Cantidad {
					// se descuenta del inventario la cantidad vendida
					inventario[producto].Stock -= orden.Cantidad

					// se responde con mensaje de éxito
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(Respuesta{
                        Mensaje: fmt.Sprintf("Venta exitosa. Quedan %d unidades.", inventario[producto].Stock),
                    })
					return
				} else {
					// si no hay suficiente stock, se responde con error	
                    respuesta.Mensaje = "Stock insuficiente"
                    w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(respuesta)
				}
                return
			} 
		}

		if !productoEncontrado {
			respuesta.Mensaje = "Producto no encontrado"
            w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(respuesta)
			return
		}
		
	})

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Bienvenido al Inventory Manager de Mercado Libre (Beta)")
    })

    fmt.Println("Servidor corriendo en el puerto 8080...")
    http.ListenAndServe(":8080", nil)
}