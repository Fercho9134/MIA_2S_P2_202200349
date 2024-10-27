package Analyzer

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"proyecto1/DiskManagement"
	"proyecto1/FileSystem"
	"proyecto1/Structs"
	"proyecto1/User"
	"proyecto1/Utilities"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

type CommandRequest struct {
	Commands []string `json:"commands"`
}

// Estructura para el JSON de respuesta
type CommandResponse struct {
	Command string `json:"command"`
	Message string `json:"message"`
}

//)
// AnalyzeHandler maneja la solicitud HTTP y ejecuta los comandos
func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Recibiendo solicitud")

	var request CommandRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	var particionesListadas []DiskManagement.PartitionInfo
	if err != nil {
		http.Error(w, "Error decodificando JSON", http.StatusBadRequest)
		log.Println("Error decodificando JSON:", err)
		return
	}

	var responses []CommandResponse
	var mensaje string

	for _, command := range request.Commands {
		//Antes de ejecutar el comando reviamos si esta linea es un comentario
		//Los comentarios tendrán un # al inicio
		if strings.HasPrefix(command, "#") {
			responses = append(responses, CommandResponse{
				Command: "Comentario",
				Message: fmt.Sprintf("> Comentario: %s", command),
			})
			log.Printf("Comentario: %s\n", command)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(responses)
			continue
		}

		commandName, params := getCommandAndParams(command)
		log.Println("Ejecutando comando:", commandName, "con parámetros:", params)
		mensaje = fmt.Sprintf("> Comando %s con parámetros: %s ejecutado exitosamente", commandName, params)
		particionesMontadasTxt := "\n> Particiones montadas:\n"
		particionesListadas, err = AnalyzeCommnad(commandName, params)

		//Imprimimos el comando ejecutado, parametros y error si lo hubo
		log.Printf("Comando: %s, Parámetros %s, Error: %v, commandName: %s\n", command, params, err, commandName)
		if err != nil {


			if commandName == "mount" {
				particionesMontadas := DiskManagement.GetMountedPartitions()
				for _, particiones := range particionesMontadas {
					for _, particion := range particiones {
						particionesMontadasTxt += fmt.Sprintf("Path: %s, Name: %s, ID: %s, Status: %d\n", particion.Path, particion.Name, particion.ID, particion.Status)

					}
				}
				//Devolvemos el mensaje de error y las particiones montadas
				responses = append(responses, CommandResponse{
					Command: commandName,
					Message: fmt.Sprintf("> %s\n%s", err.Error(), particionesMontadasTxt),
				})
			} else if commandName == "cat"{
				responses = append(responses, CommandResponse{
					Command: commandName,
					Message: fmt.Sprintf("> %s", err.Error()),
				})
			
			}else if commandName == "unmount" {
				particionesMontadas := DiskManagement.GetMountedPartitions()
				log.Println("Se ejecutó unmount con error")
				log.Println("Parametros", params)
				for _, particiones := range particionesMontadas {
					for _, particion := range particiones {
						particionesMontadasTxt += fmt.Sprintf("Path: %s, Name: %s, ID: %s, Status: %d\n", particion.Path, particion.Name, particion.ID, particion.Status)
					}
				}
				//Devolvemos el mensaje de error y las particiones montadas
				responses = append(responses, CommandResponse{
					Command: commandName,
					Message: fmt.Sprintf("> %s\n%s", err.Error(), particionesMontadasTxt),
				})
			} else if commandName == "login" || commandName == "logout" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"response": false,
					"message": err.Error(),
				})
				return
			} else if commandName == "listpartitions" {
				//Si particionesListadas es nil, entonces hubo un error
				if particionesListadas == nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"response": false,
					})
					return
				} else{
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"response": particionesListadas,
					})
					return
				}
			}else {
				responses = append(responses, CommandResponse{
					Command: commandName,
					Message: fmt.Sprintf("> %s", err.Error()),
				})
			}


		} else {
			if commandName == "mount" {
				particionesMontadas := DiskManagement.GetMountedPartitions()

				for _, particiones := range particionesMontadas {
					for _, particion := range particiones {
						particionesMontadasTxt += fmt.Sprintf("\tPath: %s, Name: %s, ID: %s, Status: %d\n", particion.Path, particion.Name, particion.ID, particion.Status)
					}
				}

				responses = append(responses, CommandResponse{
					Command: commandName,
					Message: fmt.Sprintf("%s\n%s", mensaje, particionesMontadasTxt),
				})

			} else if commandName == "unmount" {
				particionesMontadas := DiskManagement.GetMountedPartitions()
				log.Println("Se ejecutó unmount con exito")
				log.Println("Parametros", params)


				for _, particiones := range particionesMontadas {
					for _, particion := range particiones {
						particionesMontadasTxt += fmt.Sprintf("\tPath: %s, Name: %s, ID: %s, Status: %d\n", particion.Path, particion.Name, particion.ID, particion.Status)
					}
				}

				responses = append(responses, CommandResponse{
					Command: commandName,
					Message: fmt.Sprintf("%s\n%s", mensaje, particionesMontadasTxt),
				})
			} else if commandName == "login" || commandName == "logout" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"response": true,
					"message": "Inicio de sesión exitoso",
				})
				return
			} else if commandName == "listpartitions" {
				log.Println("Listando particiones")
				log.Println("Particiones listadas", particionesListadas)
				log.Println("params", params)
				particionesListadas = fn_listPartitions(params)
				log.Println("Particiones listadas", particionesListadas)
				//Si particionesListadas es nil, entonces hubo un error
				if particionesListadas == nil {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"response": false,
					})
					return
				} else{
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]interface{}{
						"response": particionesListadas,
					})
					return
				}
			}else {
				responses = append(responses, CommandResponse{
					Command: commandName,
					Message: mensaje,
				})
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responses)
	}
}

func ImprimirHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hola mundo")
}

func getCommandAndParams(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) > 0 {
		command := strings.ToLower(parts[0])
		params := strings.Join(parts[1:], " ")
		return command, params
	}
	return "", input
}

func AnalyzeCommnad(command string, params string)  ([]DiskManagement.PartitionInfo, error){
	if strings.Contains(command, "mkdisk") {
		return nil, fn_mkdisk(params)
	} else if strings.Contains(command, "rmdisk") {
		return nil, fn_rmdisk(params)
	} else if strings.Contains(command, "unmount") {
		return nil, fn_unmount(params)
	}else if strings.Contains(command, "fdisk") {
		return nil, fn_fdisk(params)
	} else if strings.Contains(command, "mount") {
		return nil, fn_mount(params)
	} else if strings.Contains(command, "mkfs") {
		return nil, fn_mkfs(params)
	} else if strings.Contains(command, "login"){
		return nil, fn_login(params)
	} else if strings.Contains (command, "logout"){
		return nil, User.Logout()
	}else if strings.Contains(command, "cat") {
		return nil, fnCat(params)
	}else if strings.Contains(command, "rep") {
		return nil, fn_rep(params)
	}else if strings.Contains(command, "readmbr"){
		return nil, fn_readmbr(params)
	}else if strings.Contains(command, "mkusr"){
		return nil, fn_mkusr(params)
	} else if strings.Contains(command, "listpartitions") {
		return fn_listPartitions(params), nil
	}else {
		return nil, fmt.Errorf("Error: Comando %s inválido o no encontrado", command)
	}
}

func fn_listPartitions(params string) []DiskManagement.PartitionInfo{
	//Debemos obtener el path del disco para listar todas sus particiones y devolverlas en fomato JSON
	fs := flag.NewFlagSet("listPartitions", flag.ExitOnError)
	path := fs.String("path", "", "Ruta del disco")

	matches := re.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := strings.ToLower(match[1])
		flagValue := strings.ToLower(match[2])
		flagValue = strings.Trim(flagValue, "\"")
		switch flagName {
		case "path":
			fs.Set(flagName, flagValue)
		default:
			return nil
		}
	}

	if *path == "" {
		return nil
	}

	//Obtenemos las particiones del disco
	particiones := DiskManagement.ListPartitions(*path)

	return particiones
}

func fn_mkusr(input string) error{
	// Definir flags
	fs := flag.NewFlagSet("mkusr", flag.ExitOnError)
	user := fs.String("user", "", "Nombre del usuario a crear")
	pass := fs.String("pass", "", "Contraseña del usuario")
	grp := fs.String("grp", "", "Grupo del usuario")

	// Parsear el input
	matches := re.FindAllStringSubmatch(input, -1)

	// Procesar el input para asignar valores a las flags
	for _, match := range matches {
		flagName := strings.ToLower(match[1])
		flagValue := strings.ToLower(match[2])
		flagValue = strings.Trim(match[2], "\"") // Limpiar el input de comillas

		switch flagName {
		case "user", "pass", "grp":
			fs.Set(flagName, flagValue)
		default:
			return fmt.Errorf("Error: Flag %s no encontrada", flagName)
		}
	}

	// Validar que las flags requeridas no estén vacías
	if *user == "" || *pass == "" || *grp == "" {
		return fmt.Errorf("Error: Los campos user, pass y grp son obligatorios")
	}

	// Verificar que los campos no excedan los 10 caracteres
	if len(*user) > 10 {
		return fmt.Errorf("Error: El usuario excede el máximo de 10 caracteres")
	}
	if len(*pass) > 10 {
		return fmt.Errorf("Error: La contraseña excede el máximo de 10 caracteres")
	}
	if len(*grp) > 10 {
		return fmt.Errorf("Error: El grupo excede el máximo de 10 caracteres")
	}

	// Crear el nuevo usuario en el formato correcto
	newUser := fmt.Sprintf("2,U,%s,%s,%s", *grp, *user, *pass)

	// Llamar a la función que maneja la creación del usuario
	err := User.MkusrCommand("/users.txt", newUser)
	if err != nil {
		return fmt.Errorf("Error: %s", err.Error())
	}

	return nil
}



func fnCat(params string) error {
    // Expresión regular para reconocer cualquier -fileN="path"
    reC := regexp.MustCompile(`-file(\d+)=("[^"]+"|\S+)`)
    matches := reC.FindAllStringSubmatch(params, -1)

    var files []string

    for _, match := range matches {
        flagValue := strings.Trim(match[2], "\"") // Quitamos las comillas si las tiene
        files = append(files, flagValue) // Añadimos el archivo a la lista
    }

    if len(files) == 0 {
        return fmt.Errorf("Error: No se encontraron archivos para mostrar")
    }

    // Aquí se llama a la función User.Cat que procesa la lista de archivos
    err := User.Cat(files)
    if err != nil {
        return fmt.Errorf("%s", err.Error())
    }

    return nil
}

func fn_login(params string) error {
	fs := flag.NewFlagSet("login", flag.ExitOnError)
	user := fs.String("user", "", "Usuario")
	pass := fs.String("pass", "", "Contraseña")
	id := fs.String("id", "", "ID")

	matches := re.FindAllStringSubmatch(params, -1)

	for _, match := range matches {
		flagName := strings.ToLower(match[1])
		flagValue := strings.ToLower(match[2])
		flagValue = strings.Trim(flagValue, "\"")
		switch flagName {
		case "user", "pass", "id":
			fs.Set(flagName, flagValue)
		default:
			return fmt.Errorf("Error: Flag %s no encontrada", flagName)
		}
	}

	if *user == "" {
		return fmt.Errorf("Error: User es obligatorio")
	}
	if *pass == "" {
		return fmt.Errorf("Error: Password es obligatorio")
	}
	if *id == "" {
		return fmt.Errorf("Error: ID es obligatorio")
	}

	err := User.Login(*user, *pass, *id)
	if err != nil {
		return fmt.Errorf("Error: %s", err.Error())
	}

	return nil
}

func fn_unmount(params string) error {
	fs := flag.NewFlagSet("unmount", flag.ExitOnError)
	id := fs.String("id", "", "ID particion a desmontar")
	matches := re.FindAllStringSubmatch(params, -1)

	//print
	fmt.Println("Matches", matches)

	for _, match := range matches {
		flagName := strings.ToLower(match[1])
		flagValue := strings.ToLower(match[2])
		flagValue = strings.Trim(flagValue, "\"")
		switch flagName {
		case "id":
			log.Println("ID", flagValue)
			fs.Set(flagName, flagValue)
		default:
			log.Println("Error: Flag %s no encontrada", flagName)
			log.Println("Value", flagValue)
			return fmt.Errorf("Error: Flag %s no encontrada", flagName)
		}
	}

	if *id == "" {
		return fmt.Errorf("Error: ID es obligatorio")
	}

	err := DiskManagement.Unmount(*id)
	if err != nil {
		return fmt.Errorf("Error: %s", err.Error())
	}

	return nil
}

func fn_mkfs(params string) error {
	fs := flag.NewFlagSet("mkfs", flag.ExitOnError)
	id := fs.String("id", "", "ID")
	type_ := fs.String("type", "", "Tipo")
	fileS := fs.String("fs", "2fs", "FS")

	matches := re.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := strings.ToLower(match[1])
		flagValue := strings.ToLower(match[2])
		flagValue = strings.Trim(flagValue, "\"")
		switch flagName {
		case "id", "type", "fs":
			fs.Set(flagName, flagValue)
		default:
			return fmt.Errorf("Error: Flag %s no encontrada", flagName)
		}

	}

	if *id == "" {
		return fmt.Errorf("Error: ID es obligatorio")
	}

	//Type puede ser vacio o full
	if *type_ != "" && *type_ != "full" {
		return fmt.Errorf("Error: Type debe ser 'full'")
	}

	err := FileSystem.Mkfs(*id, *type_, *fileS)

	if err != nil {
		return fmt.Errorf("Error: %s", err.Error())
	}
	return nil
}

func fn_mkdisk(params string) error {
	// Definir flag
	fs := flag.NewFlagSet("mkdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamaño")
	fit := fs.String("fit", "ff", "Ajuste")
	unit := fs.String("unit", "m", "Unidad")
	path := fs.String("path", "", "Ruta")

	// Encontrar la flag en el input
	matches := re.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := strings.ToLower(match[1])
		flagValue := strings.ToLower(match[2])
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size", "fit", "unit", "path":
			fs.Set(flagName, flagValue)
		default:
			return fmt.Errorf("Error: Flag %s no encontrada", flagName)
		}
	}

	// Validaciones
	if *size <= 0 {
		return fmt.Errorf("Error: La cantidad debe ser mayor a 0")
	}
	if *fit != "bf" && *fit != "ff" && *fit != "wf" {
		return fmt.Errorf("Error: El fit debe ser 'bf', 'ff', o 'wf'")
	}
	if *unit != "k" && *unit != "m" {
		return fmt.Errorf("Error: Las unidades deben ser 'k' o 'm'")
	}
	if *path == "" {
		return fmt.Errorf("Error: La ruta es obligatoria")
	}

	// Llamar a la función
	err:= DiskManagement.Mkdisk(*size, *fit, *unit, *path)
	if err != nil {
		return fmt.Errorf("Error: %s", err.Error())
	}
	return nil
}

func fn_rmdisk(params string) error {
	fs := flag.NewFlagSet("rmdisk", flag.ExitOnError)
	path := fs.String("path", "", "Ruta")

	matches := re.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := strings.ToLower(match[1])
		flagValue := strings.ToLower(match[2])
		flagValue = strings.Trim(flagValue, "\"")
		switch flagName {
		case "path":
			fs.Set(flagName, flagValue)
		default:
			return fmt.Errorf("Error: Flag %s no encontrada", flagName)
		}
	}

	if *path == "" {
		return fmt.Errorf("Error: La ruta es obligatoria")
	}

	err := DiskManagement.Rmdisk(*path)
	if err != nil {
		return fmt.Errorf("Error: %s", err.Error())
	}
	return nil
}

func fn_fdisk(params string) error {
	// Definir flags
	fs := flag.NewFlagSet("fdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamaño")
	path := fs.String("path", "", "Ruta")
	name := fs.String("name", "", "Nombre")
	unit := fs.String("unit", "k", "Unidad")
	type_ := fs.String("type", "p", "Tipo")
	fit := fs.String("fit", "", "Ajuste")
	delete_ := fs.String("delete", "", "Eliminar")
	add := fs.Int("add", 0, "Agregar")

	// Encontrar los flags en el input
	matches := re.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := strings.ToLower(match[1])
		flagValue := strings.ToLower(match[2])
		flagValue = strings.Trim(flagValue, "\"")
		switch flagName {
		case "size", "path", "name", "unit", "type", "fit", "delete", "add":
			fs.Set(flagName, flagValue)
		default:
			return fmt.Errorf("Error: Flag %s no encontrada", flagName)
		}
	}

	if *delete_ != "" {
		if *path == "" || *name == "" {
			return fmt.Errorf("Error: Path y Name son obligatorios para eliminar")
		}
		err := DiskManagement.DeletePartition(*path, *name, *delete_)
		if err != nil {
			return fmt.Errorf("Error: %s", err.Error())
		}
		return nil
	}

	if *add != 0 {
		if *path == "" || *name == "" {
			return fmt.Errorf("Error: Path y Name son obligatorios para agregar")
		}
		err := DiskManagement.ModifyPartition(*path, *name, *add, *unit)
		if err != nil {
			return fmt.Errorf("Error: %s", err.Error())
		}
		return nil
	}
	// Validaciones
	if *size <= 0 {
		return fmt.Errorf("Error: Size debe ser mayor a 0")
	}
	if *path == "" {
		return fmt.Errorf("Error: Path es obligatorio")
	}
	if *fit == "" {
		*fit = "wf"
	}
	if *fit != "bf" && *fit != "ff" && *fit != "wf" {
		return fmt.Errorf("Error: Fit debe ser 'bf', 'ff', o 'wf'")
	}
	if *unit != "k" && *unit != "m" && *unit != "b" {
		return fmt.Errorf("Error: Unidad debe ser 'k', 'm', o 'b'")
	}
	if *type_ != "p" && *type_ != "e" && *type_ != "l" {
		return fmt.Errorf("Error: Tipo debe ser 'p', 'e', o 'l'")
	}
	if *name == "" {
		return fmt.Errorf("Error: Name es obligatorio")
	}

	// Llamar a la función
	err := DiskManagement.Fdisk(*size, *path, *name, *unit, *type_, *fit)
	if err != nil {
		return fmt.Errorf("Error: %s", err.Error())
	}

	return nil
}

func fn_readmbr(params string) error {
	fs := flag.NewFlagSet("readmbr", flag.ExitOnError)
	path := fs.String("path", "", "Ruta")

	matches := re.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := strings.ToLower(match[1])
		flagValue := strings.ToLower(match[2])
		flagValue = strings.Trim(flagValue, "\"")
		switch flagName {
		case "path":
			fs.Set(flagName, flagValue)
		default:
			return fmt.Errorf("Error: Flag %s no encontrada", flagName)
		}
	}

	if *path == "" {
		return fmt.Errorf("Error: Path es obligatorio")
	}

	err := DiskManagement.ReadMBR(*path)
	if err != nil {
		return fmt.Errorf("Error: %s", err.Error())
	}

	return nil
}

func fn_mount(params string) error {
	fs := flag.NewFlagSet("mount", flag.ExitOnError)
	path := fs.String("path", "", "Ruta")
	name := fs.String("name", "", "Nombre de la partición")

	fs.Parse(os.Args[1:])
	matches := re.FindAllStringSubmatch(params, -1)

	for _, match := range matches {
		flagName := strings.ToLower(match[1])
		flagValue := strings.ToLower(match[2])
		flagValue = strings.Trim(flagValue, "\"")
		switch flagName {
		case "path", "name":
			fs.Set(flagName, flagValue)
		default:
			return fmt.Errorf("Error: Flag %s no encontrada", flagName)
		}
	}

	if *path == "" {
		return fmt.Errorf("Error: Path es obligatorio")
	}

	if *name == "" {
		return fmt.Errorf("Error: Name es obligatorio")
	}

	// Convertir el nombre a minúsculas antes de pasarlo al Mount
	nombreMinuscula := strings.ToLower(*name)
	err := DiskManagement.Mount(*path, nombreMinuscula)
	if err != nil {
		return fmt.Errorf("Error: %s", err.Error())
	}
	return nil
}

func fn_rep(params string) error {
	fs := flag.NewFlagSet("rep", flag.ExitOnError)
	name := fs.String("name", "", "Nombre")
	path := fs.String("path", "", "Ruta")
	id := fs.String("id", "", "ID")
	path_file_ls := fs.String("path_file_ls", "", "Ruta del archivo")

	matches := re.FindAllStringSubmatch(params, -1)
	for _, match := range matches {
		flagName := strings.ToLower(match[1])
		flagValue := strings.ToLower(match[2])
		flagValue = strings.Trim(flagValue, "\"")
		switch flagName {
		case "name", "path", "id", "path_file_ls":
			fs.Set(flagName, flagValue)
		default:
			return fmt.Errorf("Error: Flag %s no encontrada", flagName)
		}
	}

	if *name == "" {
		return fmt.Errorf("Error: Name es obligatorio")
	}
	if *path == "" {
		return fmt.Errorf("Error: Path es obligatorio")
	}
	if *id == "" {
		return fmt.Errorf("Error: ID es obligatorio")
	}

	//Verificamos si la particion con la id dada esta montada
	montada := false
	var pathDisco string
	for _, particiones := range DiskManagement.GetMountedPartitions() {
		for _, particion := range particiones {
			if particion.ID == *id {
				montada = true
				pathDisco = particion.Path
			}
		}
	}

	if !montada {
		return fmt.Errorf("Error: La partición con ID %s no está montada", *id)
	}

	reportsDir := filepath.Dir(*path)
	err := os.MkdirAll(reportsDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Error: %s", err.Error())
	}

	switch *name {
	case "mbr":
		file, err := Utilities.OpenFile(pathDisco)
		if err != nil {
			return fmt.Errorf("Error: %s", err.Error())
		}
		defer file.Close()

		var TempMBR Structs.MRB
		if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
			return fmt.Errorf("Error: %s", err.Error())
		}

		var ebrs []Structs.EBR
		for i := 0; i < 4; i++ {
			if string(TempMBR.Partitions[i].Type[:]) == "e" {
				log.Println("Partición extendida encontrada", string(TempMBR.Partitions[i].Name[:]))

				ebrPosition := TempMBR.Partitions[i].Start
				ebrCounter := 1

				//Leemos todos los ebrs de la partición extendida
				for ebrPosition != -1 {
					log.Println("Leyendo EBR en posicion", ebrPosition)
					var TempEBR Structs.EBR
					if err := Utilities.ReadObject(file, &TempEBR, int64(ebrPosition)); err != nil {
						return fmt.Errorf("Error: %s", err.Error())
					}

					ebrs = append(ebrs, TempEBR)
					Structs.PrintEBR(TempEBR)

					ebrPosition = TempEBR.PartNext
					ebrCounter++

					if ebrPosition == -1 {
						break
					}
				}
			}

		}

		pathReporte := *path
		if err := Utilities.GenerateReportMBR(TempMBR, ebrs, pathReporte, file); err != nil {
			return fmt.Errorf("Error: %s", err.Error())
		} else {
			log.Println("Reporte MBR generado exitosamente")
			dotFile := strings.TrimSuffix(pathReporte, filepath.Ext(pathReporte)) + ".dot"
			outupPng := strings.TrimSuffix(pathReporte, filepath.Ext(pathReporte)) + ".png"

			cmd := exec.Command("dot", "-Tpng", dotFile, "-o", outupPng)
			err := cmd.Run()
			if err != nil {
				return fmt.Errorf("Error: %s", err.Error())
			} else {
				log.Println("Imagen generada exitosamente")
			}
		}

	case "disk":
		//Generamos el reporte del disco
		file, err := Utilities.OpenFile(pathDisco)

		fileName := filepath.Base(pathDisco)

		if err != nil {
			return fmt.Errorf("Error: %s", err.Error())
		}
		defer file.Close()

		var TempMBR Structs.MRB
		if err := Utilities.ReadObject(file, &TempMBR, 0); err != nil {
			return fmt.Errorf("Error: %s", err.Error())
		}

		var ebrs []Structs.EBR
		for i := 0; i < 4; i++ {
			if string(TempMBR.Partitions[i].Type[:]) == "e" {
				log.Println("Partición extendida encontrada", string(TempMBR.Partitions[i].Name[:]))

				ebrPosition := TempMBR.Partitions[i].Start

				//Leemos todos los ebrs de la partición extendida
				for ebrPosition != -1 {
					log.Println("Leyendo EBR en posicion", ebrPosition)
					var TempEBR Structs.EBR
					if err := Utilities.ReadObject(file, &TempEBR, int64(ebrPosition)); err != nil {
						return fmt.Errorf("Error: %s", err.Error())
					}

					ebrs = append(ebrs, TempEBR)

					ebrPosition = TempEBR.PartNext

					if ebrPosition == -1 {
						break
					}
				}
			}

		}

		totalDiskSize := TempMBR.MbrSize
		pathReporte := *path
		if err := Utilities.GenerateReportDisk(TempMBR, ebrs, pathReporte, file, totalDiskSize, fileName); err != nil {
			return fmt.Errorf("Error: %s", err.Error())
		} else {
			log.Println("Reporte Disk generado exitosamente")
			dotFile := strings.TrimSuffix(pathReporte, filepath.Ext(pathReporte)) + ".dot"
			outupPng := strings.TrimSuffix(pathReporte, filepath.Ext(pathReporte)) + ".png"

			cmd := exec.Command("dot", "-Tpng", dotFile, "-o", outupPng)
			err := cmd.Run()
			if err != nil {
				return fmt.Errorf("Error: %s", err.Error())
			} else {
				log.Println("Imagen generada exitosamente")
			}
		}

		log.Println(path_file_ls)

	default:
		return fmt.Errorf("Error: Reporte %s no encontrado", *name)
	}

	return nil
}
