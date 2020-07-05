package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

type casos struct {
	Name          string `json:"Nombre"`
	Depto         string `json:"Departamento"`
	Edad          int    `json:"Edad"`
	FormaContagio string `json:"Forma_de_contagio"`
	Estado        string `json:"Estado"`
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./public"))
	//fmt.Fprintf(w, "Welcome home!")
	//http.Handle("/", http.FileServer(http.Dir("./public")))

}

var url string
var wg sync.WaitGroup

func main() {
	var rutaArchivo string
	var noHilos, cantCasos int

	fmt.Printf("URL Balanceador de carga: ")
	fmt.Scanf("%s\n", &url)
	fmt.Printf("Cantidad de hilos: ")
	fmt.Scanf("%d\n", &noHilos)
	fmt.Printf("Cantidad de solicitudes a enviar: ")
	fmt.Scanf("%d\n", &cantCasos)
	fmt.Printf("Ruta del archivo: ")
	fmt.Scanf("%s\n", &rutaArchivo)

	leerArchivo(url, noHilos, cantCasos, rutaArchivo)

}
func leerArchivo(url string, noHilos int, cantCasos int, rutaArchivo string) {
	casitos := getCasos(rutaArchivo, cantCasos)
	fmt.Println(casitos)
	aEnviar := cantCasos / noHilos
	fmt.Println(string(aEnviar))
	crearHilos(casitos, aEnviar, noHilos, cantCasos)

}
func crearHilos(caso []casos, aEnviar int, noHilos int, totalCasos int) {
	fmt.Println(string(noHilos))
	fmt.Println(string(aEnviar))

	if noHilos < 2 {
		wg.Add(noHilos)
		fmt.Println("Solo un hilo")
		enviar(caso, 0, totalCasos)
		wg.Wait()
	} else {
		fmt.Println("Mas de un hilo: ", string(noHilos))
		wg.Add(noHilos)
		for i := 0; i < noHilos; i++ {

			if (i + 1) == noHilos {
				go enviar(caso, aEnviar*i, totalCasos) //enviar ultimos
			} else if i == 0 {
				enviar(caso, i, aEnviar) //enviar primero
			} else {
				go enviar(caso, aEnviar*i, aEnviar*(i+1)) //enviar intermedios
			}
		}
		wg.Wait()

	}
}

//funcion que recibe desde donde hasta donde de los datos se van a enviar
//en cada hilo a ejecutar
func enviar(caso []casos, desde int, cantEnvio int) {
	for i := desde; i < cantEnvio; i++ {
		datosJson, _ := json.Marshal(caso[i])
		fmt.Println("%s\n", string(datosJson))
		//se envian datos como json al url solicitado
		_, err := http.Post(url, "application/json", bytes.NewBuffer(datosJson))
		if err != nil {
			fmt.Printf("Error al enviar: %s\n", err)
		}
		time.Sleep(time.Millisecond * 10)
	}
	wg.Done()
}

func (c casos) toString() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(bytes)
}

func getCasos(ruta string, cantCasos int) []casos {
	fmt.Println("funcionGetCasos")
	caso := make([]casos, cantCasos)
	raw, err := ioutil.ReadFile(ruta)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	json.Unmarshal(raw, &caso)

	aux := len(caso)
	j := 0

	for i := len(caso); i < cantCasos; i++ {
		caso = append(caso, caso[j])
		j++

		if j == aux {
			j = 0
		}
	}

	return caso
}

/*
func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newEvent)
	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newEvent)
}
*/
