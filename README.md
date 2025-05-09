Ejercicio Final – API de Ventas en Go
Partimos del proyecto base (go_parte_3) que ya expone un CRUD de usuarios usando Gin. Ahora vamos a extender nuestro sistema para manejar ventas (sales/transactions).
Objetivos
Diseñar e implementar una nueva API sales-api en Go. La misma deberá brindar los siguientes endpoints:

Creación de una nueva venta: POST /sales

Recibe JSON con user_id y amount.

Valida que el user_id exista llamando a GET /users/:id.

Valida que el monto no sea 0.

Asigna aleatoriamente uno de estos estados: pending, approved, rejected.

Debe generar el ID automáticamente.

Guarda la venta usando una capa de storage (puede ser en memoria).

Devuelve un json parecido a esto:



El campo created_at y updated_at deberán ser del tipo time.Time.

El campo version deberá ser un int. Empieza en 1.

Actualización de una venta: PATCH /sales/:id

Lo único que debe recibir este endpoint es un JSON con un campo status el cual sólo podrá afectar a las ventas que queden en estado pending. 

Las únicas transiciones posibles serán pending -> approved o pending -> rejected.

Deberá controlar que la venta exista y que además esté en el estado correcto para hacer el pasaje.

La respuesta de este endpoint deberá ser la misma que en el punto 1, sólo que con la fecha de actualización, versión y estado actualizado.

Search de ventas: GET /sales?user_id={{user_id}}&status={{status}}

Este endpoint deberá buscar todas las ventas de un usuario particular.

En caso de que el campo status se envíe, entonces se deberá realizar un filtro, caso contrario deberá traer todas las ventas en cualquier estado.

En caso de que se pase un estado para filtrar, deberá validar que este es correcto.

En la respuesta de este endpoint, se deberá incluir un atributo metadata en donde se debe incluir:
Cantidad de ventas: quantity.
Cantidad de ventas aprobadas: approved.
Cantidad de ventas rechazadas: rejected.
Cantidad de ventas pendientes: pending.
Monto total (suma del monto de todas las ventas): total_amount.

La firma de respuesta de este endpoint deberá ser de la siguiente manera:



Detalles Técnicos
El código de la API se deberá subir a un repositorio público de Github de alguno de los miembros del grupo de trabajo para que pueda ser revisado.

Para crear un cliente HTTP (para poder comunicarte en código con la API de usuarios) hacer uso de la librería: https://github.com/go-resty/resty.

Usar la librería https://github.com/google/uuid para crear los UUIDs de las ventas.

No se exigirá ningún tipo de estructura particular. No obstante, pueden utilizar la API de usuarios como modelo.

El guardado de las entidades “ventas” puede ser en memoria, tal cual es en la API de usuarios.

El proyecto tendrá que tener AL MENOS un test unitario y un test de integración.

Un test unitario intentando crear una venta con user inexistente.
Un test de integración que levante Gin y pruebe un flujo completo de POST → PATCH → GET en el happy path.

Se apreciará el uso de técnicas de observabilidad vistas en la última clase (logs). Pueden utilizar el logger de Uber https://github.com/uber-go/zap.

Prestar especial atención a las respuestas HTTP ante cada error en cada uno de los endpoints. El detalle es el siguiente:

POST /sales

Si el usuario no existe: 
400 bad request.

Si el monto es inválido o cualquier otro error de deserialización:
400 bad request.

Cualquier otro error que se pueda producir: 
500 internal server error.

Se creó satisfactoriamente:
201 status created.

PATCH /sales/:id

Si la venta no existe: 
404 not found.

Si el estado a transicionar es inválido o hay un error en el body:
400 bad request.

Si se quiere realizar una transición inválida (por ejemplo pasar de approved a rejected): 
409 status conflict.

Cualquier otro error que se pueda producir: 
500 internal server error.

Se actualizó correctamente:
200 status OK.

GET /sales?user_id={{user_id}}&status={{status}}

Si el estado para buscar es inválido.
400 bad request.

Cualquier otro error que se pueda producir: 
500 internal server error.

Traiga o NO traiga resultados:
200 status OK.

En caso de que no traiga resultados entonces el atributo results deberá mostrar un array vacío [].
