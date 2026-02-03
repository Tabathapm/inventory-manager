# 📦 Inventory Manager (PostgreSQL Monolith)

Segunda iteración del proyecto Inventory Manager. En esta etapa, introdujimos persistencia de datos real utilizando **PostgreSQL** y contenedorización básica.

## ⚙️ Características Técnicas

* **Arquitectura:** Monolito (Todo el código en un solo servicio).
* **Base de Datos:** PostgreSQL 15 (Imagen Alpine).
* **Infraestructura:** `docker-compose.yml` para vincular la App con la DB.
* **Driver:** `lib/pq` para conectar Go con SQL.

## 🛠️ Instalación y Ejecución

No necesitas instalar Go ni Postgres localmente, solo Docker.

1.  **Ejecutar el entorno:**
    ```bash
    docker compose up --build
    ```

2.  **Verificación:**
    El sistema estará disponible en `http://localhost:8080`.
    La base de datos corre en el puerto `5432` (user: `postgres`, pass: `secret`).

## 📡 Endpoints Disponibles

| Método | Endpoint | Descripción |
| :--- | :--- | :--- |
| `POST` | `/productos` | Crea un nuevo producto en la DB. |
| `GET` | `/productos` | Lista todos los productos guardados. |
| `POST` | `/vender` | Registra una venta y actualiza el stock (SQL Transaction). |

### Ejemplo de JSON para Venta:
```json
{
    "producto_id": 1,
    "cantidad": 5
}
```