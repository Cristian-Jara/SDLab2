# SDLab2

### Integrantes

- Cristian Jara 201704563-9

- Sebastian Muñoz 201473503-0

- Maria Riveros 201704585-k


### Instrucciones para funcionamiento

Primero se debe conectar el lider en la máquina virtual dist86 con ip 10.6.40.227 con el comando "make runLider", con esto ya se puede levantar los demás componentes sin ningún orden en específico (Solo procurar dejar los jugadores para el final). De todas maneras se detallará un posible orden por máquina:

1. Conectar Lider en la máquina virtual dist86 con "make runLider", luego en la misma máquina conectar un DataNode con "make runDN"
2. Para la máquina dist85 se debe levantar el pozo con "make runPozo" y un DataNode con "make runDN"
3. En la máquina dist87 se debe conectar el NameNode con el comando "make runNN"
4. Para la máquina dis88 se debe conectar el último DataNode con el comando "make runDN"
5. Ya se pueden conectar los jugadores en cualquiera de las máquinas con "make runPlayer" (Manejables)
6. Para rellenar con bots se deben crear los bots uno por uno con "make runBot" (Automáticos)

Como se puede notar el juego fue pensado para dar la posibilidad de conectar la cantidad
requerida de jugadores uno por uno (dejando la posibilidad de tener los bots y jugadores reales que se quieran), principalmente porque se habla de procesos que deben finalizar una vez pierdan (asumido).

Se procedió a abrir puertos para ver el código en funcionamiento, pero para el pozo, la creacion del usuario junto con la autentificación, en específico, la configuracion de rabbit no nos alcanzó el tiempo para realizarla. De todas formas, cabe destacar que el pozo esta implementado y funcional a nivel local. Si se sabe como configurar la máquina para que funcioné RabbitMQ se podrá probar, para ello bastaría con borrar los comentarios que tiene la funcion JugadorMuerto en "Lider/Lider.go" y realizaría las conexiones.