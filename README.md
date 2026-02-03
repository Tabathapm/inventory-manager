# 📦 Inventory Manager (Microservices Architecture)

Sistema de gestión de inventario distribuido implementado en **Go**. Esta versión representa la evolución final del proyecto, migrando de un monolito a una arquitectura de **Microservicios**.

## 🚀 Arquitectura del Sistema

El sistema está compuesto por contenedores Docker orquestados:

1.  **Inventory Service (`:8080`):**
    * Servicio principal (Backend/API Gateway).
    * Gestiona productos y orquesta la venta.
    * Se comunica con la base de datos y el servicio de pagos.
2.  **Payment Service (`:8081`):**
    * Microservicio independiente simulado.
    * Procesa transacciones y aprueba/rechaza pagos con latencia simulada.
3.  **Database (PostgreSQL):**
    * Persistencia de datos para el inventario y las órdenes.

## 🛠️ Requisitos Previos

* [Docker](https://www.docker.com/)
* [Docker Compose](https://docs.docker.com/compose/)

## ⚡ Guía de Inicio Rápido

1.  **Clonar y acceder a la rama:**
    ```bash
    git checkout feature/microservices
    ```

2.  **Levantar la infraestructura:**
    ```bash
    docker compose up --build
    ```
    _Espera a ver los logs de "Servicio de Pagos corriendo..." y "Servidor corriendo en el puerto 8080"._

## 🧪 Cómo Probar (API Reference)

### 1. Crear un Producto (Setup inicial)
**POST** `http://localhost:8080/productos`
```json
{
    "nombre": "Laptop Gamer",
    "precio": 1500.00,
    "stock": 10
}
```

### 2. Realizar una Venta (Flujo Completo)
Esta petición activa la comunicación entre microservicios.

**POST** `http://localhost:8080/vender`

```json
{
    "producto_id": 1,
    "cantidad": 1
}
```

**Respuesta esperada:**
* Si el pago es aprobado (200 OK), recibirás la orden creada.
* Si el pago falla, recibirás un error y no se descontará stock.

## 📝 Nota
Desarrollado como proyecto educativo de arquitectura de software.