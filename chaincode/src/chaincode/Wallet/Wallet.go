package main

/**

Autor: Andres Martinez Melgar

Chaincode que simula una cartera.
Se podra realizar una de las siguientes acciones:
	CrearCartera
	BorrarCartera
	EnviarDinero
	Query
	QueryOnTime(trazabilidad)

*/
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type CarteraChaincode struct {
}

// Informacion del usuario que se va a guardar en el ledger
type Cartera struct {
	ValorActual     string `json:"ValorActual"`
	ValorMaximo     string `json:"ValorMaximo"`
	Contrasena      string `json:"Contrasena"`
	FechaCreacion   string `json:"FechaCreacion"`
	Nombre          string `json:"Nombre"`
	Apellido1       string `json:"Apellido1"`
	Apellido2       string `json:"Apellido2"`
	FechaNacimiento string `json:"FechaNacimiento"`
}

// funcion principar que es llamada cuando el chaincode se instala o actualiza
func (c *CarteraChaincode) Init(stub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

//funcion que se llama siempre que se quiere reaizar una accion con el chaincode
func (c *CarteraChaincode) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	function, parameters := stub.GetFunctionAndParameters()

	if function == "query" {
		return c.query(stub, parameters)
	} else if function == "crearCartera" {
		return c.crearCartera(stub, parameters)
	} else if function == "borrarCartera" {
		return c.borrarCartera(stub, parameters)
	} else if function == "queryOnTime" {
		return c.queryOnTime(stub, parameters)
	} else if function == "enviarDinero" {
		return c.enviarDinero(stub, parameters)
	}

	return shim.Error("Funcion introducida no encontrada")
}

//funcion para enviar dinero
func (c *CarteraChaincode) enviarDinero(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	From := args[0]                            //usuario que manda dinero
	To := args[1]                              //usuario que recibe dinero
	cantidadDinero, _ := strconv.Atoi(args[2]) //cantidad de dinero enviado
	FromAsBytes, err := stub.GetState(From)
	//empiezan las comprobaciones
	if err != nil {
		return shim.Error("Fallo al recuperar  el usuario")
	}
	if FromAsBytes == nil {
		return shim.Error("El usuario no existe")
	}

	ToAsBytes, err := stub.GetState(To)
	if err != nil {
		return shim.Error("Fallo al recuperar  el usuario")
	}
	if ToAsBytes == nil {
		return shim.Error("El usuario no existe")
	}
	//una vez que se que los usuarios existen los parseo de []byte a Usuario
	usuarioFROM := Cartera{}
	usuarioTO := Cartera{}

	err = json.Unmarshal(FromAsBytes, &usuarioFROM)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = json.Unmarshal(ToAsBytes, &usuarioTO)
	if err != nil {
		return shim.Error(err.Error())
	}
	if usuarioFROM.Contrasena != args[3] {
		return shim.Error("Las contrasenas no coinciden")
	}

	dineroActualFROM, _ := strconv.Atoi(usuarioFROM.ValorActual)
	dineroActualTO, _ := strconv.Atoi(usuarioTO.ValorActual)

	if dineroActualFROM < cantidadDinero {
		return shim.Error("El usuario de origen no tiene esa cantidad de dinero")
	}

	//fin comprobaciones
	usuarioFROM.ValorActual = strconv.Itoa(dineroActualFROM - cantidadDinero)
	usuarioTO.ValorActual = strconv.Itoa(dineroActualTO + cantidadDinero)

	com1, _ := strconv.Atoi(usuarioTO.ValorActual)
	comp2, _ := strconv.Atoi(usuarioTO.ValorMaximo)
	if com1 > comp2 {
		usuarioTO.ValorMaximo = usuarioTO.ValorActual
	}
	//fin del envio de dinero

	FromAsBytes, _ = json.Marshal(usuarioFROM)
	ToAsBytes, _ = json.Marshal(usuarioTO)

	//fin del parseo a []byte
	stub.PutState(usuarioFROM.Nombre, FromAsBytes)
	stub.PutState(usuarioTO.Nombre, ToAsBytes)
	//fin del guardado de datos en el ledger
	return shim.Success(nil)
}

//funcion para crear una cartera nueva
func (c *CarteraChaincode) crearCartera(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	var i = -1
	var cartera = Cartera{
		ValorActual:     args[sumaInt(&i)],
		ValorMaximo:     args[sumaInt(&i)],
		Contrasena:      args[sumaInt(&i)],
		FechaCreacion:   args[sumaInt(&i)],
		Nombre:          args[sumaInt(&i)],
		Apellido1:       args[sumaInt(&i)],
		Apellido2:       args[sumaInt(&i)],
		FechaNacimiento: args[sumaInt(&i)]}
	carteraAsByte, _ := json.Marshal(cartera)
	stub.PutState(cartera.Nombre, carteraAsByte)
	return shim.Success(nil)
}

// funcion para borrar una cartera existente
func (c *CarteraChaincode) borrarCartera(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	//args[0] cartera a borrar
	//args[1] contrasena de seguridad
	if len(args) != 2 {
		return shim.Error("Error al borrar el usuario, se esperaba solo 2 arg")
	}
	//vemos si el usuario existe

	usuarioAsByte, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("El usuario a borrar no existe")
	}

	usuario := Cartera{}

	err = json.Unmarshal(usuarioAsByte, &usuario)
	if err != nil {
		return shim.Error(err.Error())
	}
	if usuario.Contrasena != args[1] {
		return shim.Error("Las contrasenas no coinciden")
	}
	//borramos el usuario
	stub.DelState(args[0])
	return shim.Success(nil)
}

//funcion que devuelve los ultimos datos de una cartera
func (c *CarteraChaincode) query(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	//args[0] cartera a buscar
	if len(args) != 1 {
		return shim.Error("Error en query, numero incorrecto de parametros, se esperaba solo 1.")
	}
	wallet, _ := stub.GetState(args[0])
	if wallet == nil {
		return shim.Error("No se ha encontrado ninguna coincidencia")
	}
	return shim.Success(wallet)

}

//funcion que sirve para realizar una trazabilidad de una cartera
func (c *CarteraChaincode) queryOnTime(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) < 1 {
		return shim.Error("Numero incorrecto de parametros,, se esperaba  1")
	}

	NombreWallet := args[0]

	fmt.Printf("- Trazabilidad de: %s\n", NombreWallet)
	//recupero la wallet desde el ledger
	resultsIterator, err := stub.GetHistoryForKey(NombreWallet)
	if err != nil {
		return shim.Error(err.Error())
	}
	//cierro el iterador cuando salga de la funcion
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	flag := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if flag == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Datos\":")

		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}
		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")
		buffer.WriteString(", \"Borrado?\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		flag = true
	}
	buffer.WriteString("]")

	fmt.Printf("- Trazabilidad Wallet:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
func sumaInt(i *int) int {
	*i = *i + 1
	return *i
}
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(CarteraChaincode))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
