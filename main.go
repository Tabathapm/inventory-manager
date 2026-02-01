package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "inventory-manager/database"
    "strconv"
    "strings"
    "bytes"
)

type SolicitudPago struct {
	UsuarioID int `json:"usuario_id"`
	Monto float64 `json:"monto"`
}

type RespuestaPago struct {
	Estado string `json:"estado"`
	IDTransaccion string `json:"id_transaccion"`
}

// Función "Privada" (minúscula): Solo sirve para usarla internamente
func realizarCobro(usuarioID int, monto float64) error {
	
	// Se arma la estructura con los datos que recibidos
	solicitud := SolicitudPago{
		UsuarioID: usuarioID,
		Monto:     monto,
	}

	// Convertir esa estructura a JSON (bytes)
	// json.Marshal devuelve los datos codificados y un error (si falló)
	jsonData, err := json.Marshal(solicitud)
	if err != nil {
		return fmt.Errorf("error al convertir a JSON: %v", err)
	}

	// Preparar el envío
	// Convertir los bytes en un "Reader" (un flujo de datos) que http.Post entiende
	body := bytes.NewBuffer(jsonData)
	
	// Se usa el nombre del servicio de Docker "payments" y el puerto 8081
	url := "http://payments:8081/procesar-pago"

	// Hacer la llamada (El POST real)
	// Le decimos: "Toma, acá va un JSON ('application/json') con estos datos (body)"
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		return fmt.Errorf("error de conexión con pagos: %v", err)
	}
	// "defer" asegura que se cerró la conexión al terminar la función
	defer resp.Body.Close()

	// Verificación de respuesta
	// Si el código no es 200 OK, algo salió mal (ej: tarjeta rechazada)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("el pago fue rechazado con código: %d", resp.StatusCode)
	}

	// Si se llegó hasta acá, todo ok. Se devuelve nil (sin error)
	return nil
}

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

    // --------- VENDER (Con cobro verificado)---------
	http.HandleFunc("/vender", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != "POST" {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		var orden database.OrdenDeCompra
		if err := json.NewDecoder(r.Body).Decode(&orden); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		producto, err := database.TraerProducto(orden.Producto_id)
		if err != nil {
			http.Error(w, "Producto no encontrado", http.StatusNotFound)
			return
		}

		totalACobrar := producto.Precio * float64(orden.Cantidad)

		// LLAMADA AL MICROSERVICIO DE PAGOS 
		// Se usa el ID de usuario 1 (simulado)
		fmt.Printf("Intentando cobrar $%0.2f al usuario 1...\n", totalACobrar)
		
		err = realizarCobro(1, totalACobrar)
		if err != nil {
			// Si el cobro falla: se corta todo y no se guarda la venta
			fmt.Printf("❌ Cobro fallido: %v\n", err)
			http.Error(w, "El pago fue rechazado. No se procesó la venta.", http.StatusPaymentRequired)
			return
		}

		fmt.Println("✅ Cobro exitoso. Registrando venta en base de datos...")

		// Se guarda la orden y se resta del stock.
		ordenGuardada, err := database.RegistrarOrdenDeCompra(orden)
		if err != nil {
			http.Error(w, "Error al guardar la orden: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Todo chill :D
		json.NewEncoder(w).Encode(ordenGuardada)
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