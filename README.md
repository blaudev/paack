# Paack
## Descripción
Tenemos un archivo de clientes en formato csv. Este archivo es lo suficientemente grande como para no poder cargarse completo en memoria.  

Necesitamos 2 servicios:  
El primer servicio lee este archivo y actualiza una base de datos postgres con los datos de este.  

Una vez actualizada la base de datos comunica los cambios al segundo servicio.  

El segundo servicio se encarga de enviar las incorporaciones y modificaciones a un API de un CMS. Esta API se trata de un API REST clásica que no acepta POST ni PUT masivos. Además, esta API falla de forma aleatória.  
## Resolución
### Primer servicio: PARSER
Se encarga de parsear los registros del archivo. Los inserta o actualiza en la base de datos y los enviamos al segundo servicio.  

Estos son los pasos:  

**Paso 1:**  
Se divide el archivo de clientes en n archivos de menos tamaño para poder ser tratados
sin problemas de memoria.
  
**Paso 2:**  
Se iteran todos los archivos.  
Se parsean todos los registros de cada archivo.  
Los registros se dividen en arrays de menor tamaño para enviarse como transacciones a la base de datos.  
La base de datos devuelve los regsitros actualizados indicando si se trata de registros nuevos o actualizados.  
Se utiliza concurrencia en el envío de transacciones para ganar velocidad.  
  
**Paso 3:**  
Los registros devueltos por la base de datos se dividen y se envían al segundo servicio.  
Se utiliza el protocolo RPC para mayor velocidad y seguridad.  
Se utiliza concurrencia para el envío para conseguir mayor velocidad.  

## Segundo servicio: INTEGRATOR
El servicio recibe los registros del **parser** y los envía a la API del CMS.  
Si el APi falla se repite la operación hasta un máximo de n veces.  
Si el registro es nuevo se envía un POST, si es una actualización se envía un PUT.  
Solo se repetirá un envío en caso de error por parte del API.  
  
## Ejecución
### Parser
```
$ cd ../../apps/parser
$ go build -o ../../bin/parser
$ cd ../../bin
$ ./parser -r 1000000 -f customers.csv
```
| Parámetro | Valor                            |
| --------- | -------------------------------- |
| -r        | número de registros              |
| -f        | path para el archivo de clientes |
  
**-r** indica el número de regsitros que contendrá cada archivo una vez se haya dividido el csv de clientes.    

### Integrator
Al tratarse de un servico arranca con docker-compose:
```
$ sudo docker-compose up --build
```

## Tools
Para facilitar el testing se han creado 2 tools y 2 servicios:  
### Mock
Crea un archivo de clientes para pruebas  
```
$ cd ../../tools/mock
$ go build -o ../../bin/mock
$ cd ../../bin
$ ./mock -r 1000000 -f customers.csv
```
| Parámetro | Valor                            |
| --------- | -------------------------------- |
| -r        | número de registros              |
| -f        | path para el archivo de clientes |
  
## API, POSTGRES y ADMINER
Api simula el API REST del CMS, fallando de forma random.  
Postgres y admin para testing de la base de datos.  
Arrancan con docker-compose.  
```
$ sudo docker-compose up --build
```
# Notas 
La cobertura de testing es mínima. La mayoría de testing debería ser de integración, lo cual escapa al objetivo de esta prueba.  

Se ha procurado seguir las convenciones del lenguaje. Por ejemplo, ninguna función ni struct tiene letra capital si no es estrictamente necesario.  

Cada servicio se ha desarrollado en un módulo separado por si se escala a microsevicios.  

Cada servicio se ha desarrollado en un solo package siguiendo la norma de responsabilidad única típica de Go.

No se eliminan los registros anteriores (sólo se insertan y actualizan los existentes). Creo que no era el objetivo de esta prueba.  

Los Docker no se han optimizado. Están solo para uso en desarrollo.  

Hay varias constántes que se pueden ajustar según requisitos.  
