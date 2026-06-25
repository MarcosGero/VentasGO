# API de Ventas en Go

API REST en Go para la gestión de ventas, con validación entre servicios, máquina de estados y observabilidad básica. Desarrollada como ejercicio final de un curso de Go: extiende una API de usuarios base (provista en el curso) agregando el dominio completo de ventas, su lógica de negocio y sus pruebas.

## Qué hace

El servicio expone una API de ventas (`sales-api`) que se comunica con una API de usuarios independiente para validar datos antes de operar. Cada venta atraviesa una máquina de estados con transiciones controladas, y las consultas devuelven métricas agregadas además de los resultados.

## Stack

- **Go** + **[Gin](https://github.com/gin-gonic/gin)** — framework HTTP
- **[resty](https://github.com/go-resty/resty)** — cliente HTTP para la comunicación entre servicios
- **[google/uuid](https://github.com/google/uuid)** — generación de identificadores
- **[zap](https://github.com/uber-go/zap)** — logging estructurado
- Almacenamiento **en memoria** (capa de storage desacoplada, intercambiable)

## Endpoints

### `POST /sales` — Crear una venta

Recibe `user_id` y `amount`. Valida que el usuario exista (consultando `GET /users/:id` en la API de usuarios) y que el monto sea distinto de cero. Asigna un estado inicial (`pending`, `approved` o `rejected`), genera el `id` y persiste la venta.

```json
{
  "id": "60242af5-d080-49a4-9b07-4e63992094ef",
  "user_id": "60242af5-d080-49a4-9b07-4e63992094ef",
  "amount": 100.00,
  "status": "approved",
  "created_at": "2025-04-21T20:17:29.673021-03:00",
  "updated_at": "2025-04-21T20:17:29.673022-03:00",
  "version": 1
}
```

**Respuestas:** `201` creada · `400` usuario inexistente, monto inválido o body malformado · `500` error interno.

### `PATCH /sales/:id` — Actualizar el estado de una venta

Recibe un campo `status`. Solo opera sobre ventas en estado `pending`, y únicamente permite las transiciones `pending → approved` y `pending → rejected`. Al actualizar, incrementa la versión y la fecha de modificación.

**Respuestas:** `200` actualizada · `404` venta inexistente · `409` transición inválida (ej. `approved → rejected`) · `400` body inválido · `500` error interno.

### `GET /sales?user_id={id}&status={status}` — Buscar ventas

Devuelve todas las ventas de un usuario. El filtro `status` es opcional; si se envía, se valida que sea un estado conocido. La respuesta incluye un bloque `metadata` con métricas agregadas. Si no hay resultados, `results` es un array vacío.

```json
{
  "metadata": {
    "quantity": 2,
    "approved": 1,
    "rejected": 1,
    "pending": 0,
    "total_amount": 300.0
  },
  "results": [ /* ventas */ ]
}
```

**Respuestas:** `200` con o sin resultados · `400` estado de filtro inválido · `500` error interno.

## Decisiones de diseño

- **Códigos HTTP por caso de error.** Cada endpoint distingue entre error del cliente y error del servidor, y usa `409 Conflict` específicamente para transiciones de estado inválidas en lugar de un `400` genérico, para que el consumidor de la API entienda *qué* salió mal.
- **Capa de storage desacoplada.** La persistencia está detrás de una interfaz; esta en memoria, pero se puede reemplazar por una base de datos sin tocar la lógica de negocio.
- **Comunicación entre servicios.** La validación del usuario se hace contra la API de usuarios por HTTP, simulando un escenario de microservicios en vez de asumir que el dato es válido.
- **Observabilidad.** Logging estructurado con zap para poder rastrear el flujo de las requests.

## Ejecucion

> Requiere la API de usuarios corriendo en paralelo, ya que la creación de ventas la consulta.

```bash
# 1. Levantar la API de usuarios (proyecto base)
# ← agregar acá el comando según tu estructura

# 2. Levantar la API de ventas
go run .          # ← ajustar si el entrypoint está en otra ruta
```

Por defecto la API queda escuchando en `http://localhost:8080`. ← *verificar el puerto en tu código*

Ejemplo de uso:

```bash
curl -X POST http://localhost:8080/sales \
  -H "Content-Type: application/json" \
  -d '{"user_id": "60242af5-d080-49a4-9b07-4e63992094ef", "amount": 100.0}'
```

## Pruebas

```bash
go test ./...
```

Incluye:
- **Test unitario:** intento de crear una venta con un usuario inexistente.
- **Test de integración:** ejecuta un flujo completo `POST → PATCH → GET` sobre el camino feliz.
