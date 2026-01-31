package main

import (
	"net/http" 			//para manejar el servidor web
	"encoding/json" 	//para leer y escribir JSON
	"fmt" 				//para imprimir mensajes en la consola
	"log" 				//para imprimir mensajes en la consola
	"time" 				// para manejar el tiempo
)

type SolicitudPago struct {
	UsuarioID int `json:"usuario_id"`
	Monto float64 `json:"monto"`
}

type RespuestaPago struct {
	Estado string `json:"estado"`
	IDTransaccion string `json:"id_transaccion"`
}

func main()  {
	http.HandleFunc("/procesar-pago", func(w http.ResponseWriter, r *http.Request)  {
		w.Header().Set("Content-Type", "application/json")
		
		if r.Method != "POST" {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		
		var solicitud SolicitudPago
		if err := json.NewDecoder(r.Body).Decode(&solicitud); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		fmt.Println("Cobrando...")
		time.Sleep(2*time.Second)

		var respuesta RespuestaPago
		respuesta.Estado = "aprobado"
		respuesta.IDTransaccion = "tx_12345"

		json.NewEncoder(w).Encode(respuesta)

	})

	fmt.Println("Servicio de Pagos corriendo en el servidor 8081...")
    log.Fatal(http.ListenAndServe(":8081", nil))
}