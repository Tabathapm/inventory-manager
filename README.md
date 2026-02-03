# 📦 Inventory Manager (In-Memory Version)

Versión original (MVP) del sistema de gestión de inventario.

Esta versión está diseñada para demostrar los fundamentos de **Go (Golang)** sin la complejidad de bases de datos externas ni Docker. Todos los datos se almacenan temporalmente en la memoria RAM del servidor.

## ⚠️ Aviso Importante
**Los datos son volátiles.** Si detienes o reinicias el servidor, todos los productos y ventas creados se perderán, ya que se guardan en variables `map` y `slice` en memoria.

## 🚀 Cómo Ejecutar

Requisitos: Tener [Go instalado](https://go.dev/dl/) en tu máquina.

1.  **Iniciar el servidor:**
    ```bash
    go run main.go
    ```

2.  **Probar:**
    Abre tu navegador o Postman en `http://localhost:8080/productos`.

## 📚 Conceptos Aplicados
* Estructuras (`structs`) y Tipos de Datos.
* Servidor Web básico con `net/http`.
* Manipulación de JSON (`encoding/json`).
* Lógica de negocio simple (CRUD).